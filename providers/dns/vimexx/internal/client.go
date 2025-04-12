package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const TTL = 3600 * 24

type Vimexx struct {
	clientId    string
	clientKey   string
	username    string
	password    string
	apiUrl      string
	endpoint    string
	accessToken Token
	HTTPClient  *http.Client
}

type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewClient(clientId string, clientKey string, username string, password string, apiUrl string, endpoint string) *Vimexx {
	return &Vimexx{
		clientId:  clientId,
		clientKey: clientKey,
		username:  username,
		password:  password,
		apiUrl:    apiUrl,
		endpoint:  endpoint,
	}

}

func (v *Vimexx) Login() {
	body := url.Values{}
	body.Set("grant_type", "password")
	body.Set("client_id", v.clientId)
	body.Set("client_secret", v.clientKey)
	body.Set("username", v.username)
	body.Set("password", v.password)
	body.Set("scope", "whmcs-access")

	req, _ := http.NewRequest("POST", v.apiUrl+"/auth/token", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := v.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(res.Body).Decode(&v.accessToken)
	if err != nil {
		panic(err)
	}
	//f, _ := os.Create("token.json")
	//j, _ := json.Marshal(v.accessToken)
	//f.Write(j)
	//f.Close()

}

func (v *Vimexx) LoadToken(file string) {
	f, _ := os.Open(file)
	json.NewDecoder(f).Decode(&v.accessToken)
}

type DNSRecord struct {
	Name    string `json:"name"`
	Type    string `json:"type,omitempty"`
	Content string `json:"content,omitempty"`
	Prio    string `json:"prio,omitempty"`
	TTL     int    `json:"ttl,omitempty"`
}

type DNSData struct {
	DNSRecords []DNSRecord `json:"dns_records"`
}

type DNSResponse struct {
	Message string  `json:"message"`
	Result  bool    `json:"result"`
	Data    DNSData `json:"data"`
}

func (v Vimexx) GetDNS(domain string) []DNSRecord {
	parts := strings.Split(domain, ".")
	j, _ := json.Marshal(map[string]any{
		"body": map[string]string{
			"sld": parts[0],
			"tld": parts[1],
		},
		"version": "8.6.0-release.1",
	})
	fmt.Println(string(j))
	req, _ := http.NewRequest("POST", v.apiUrl+v.endpoint+"/whmcs/domain/dns", bytes.NewReader(j))
	req.Header.Add("Authorization", "Bearer "+v.accessToken.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := v.HTTPClient.Do(req)

	if err != nil {
		panic(err)
	}
	response := DNSResponse{}
	json.NewDecoder(res.Body).Decode(&response)
	if !response.Result {
		panic("No DNS records found")
	}
	fmt.Println(response)
	return response.Data.DNSRecords
}

func (v Vimexx) SetDNS(domain string, records []DNSRecord) {
	parts := strings.Split(domain, ".")
	rr := records[:0]
	for _, record := range records {
		if record.TTL == 0 {
			record.TTL = TTL
		}
		rr = append(rr, record)
	}
	j, _ := json.Marshal(map[string]any{
		"body": map[string]any{
			"sld":         parts[0],
			"tld":         parts[1],
			"dns_records": rr,
		},
		"version": "8.6.0-release.1",
	})
	fmt.Println(string(j))
	req, _ := http.NewRequest("PUT", v.apiUrl+v.endpoint+"/whmcs/domain/dns", bytes.NewReader(j))
	req.Header.Add("Authorization", "Bearer "+v.accessToken.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := v.HTTPClient.Do(req)

	if err != nil {
		panic(err)
	}
	response := DNSResponse{}
	json.NewDecoder(res.Body).Decode(&response)
	fmt.Println(response)
	if !response.Result {
		panic("Failed Setting DNS Records")
	}
}
