// Package Vimexx implements a DNS provider for solving the DNS-01 challenge using Vimexx DNS.
package vimexx

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
	"github.com/go-acme/lego/v4/providers/dns/vimexx/internal"
)

// Environment variables names.
const (
	envNamespace = "VIMEXX_"

	EnvServerBaseURL = envNamespace + "SERVER_BASE_URL"
	EnvEndpoint      = envNamespace + "ENDPOINT"
	EnvUsername      = envNamespace + "USERNAME"
	EnvPassword      = envNamespace + "PASSWORD"
	EnvClientId      = envNamespace + "CLIENT_ID"
	EnvClientSecret  = envNamespace + "CLIENT_SECRET"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

// Config is used to configure the creation of the DNSProvider.
type Config struct {
	baseURL      string
	Endpoint     string
	Username     string
	Password     string
	ClientId     string
	ClientSecret string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	HTTPClient         *http.Client
}

// NewDefaultConfig returns a default configuration for the DNSProvider.
func NewDefaultConfig() *Config {
	return &Config{
		TTL:                env.GetOrDefaultInt(EnvTTL, 300),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, dns01.DefaultPropagationTimeout),
		PollingInterval:    env.GetOrDefaultSecond(EnvPollingInterval, dns01.DefaultPollingInterval),
		HTTPClient: &http.Client{
			Timeout: env.GetOrDefaultSecond(EnvHTTPTimeout, 30*time.Second),
		},
	}
}

// DNSProvider implements the challenge.Provider interface.
type DNSProvider struct {
	config *Config
	client *internal.Vimexx
}

// NewDNSProvider returns a DNSProvider instance configured for Vimexx.
// Credentials must be passed in the environment variables:
func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvServerBaseURL, EnvUsername, EnvPassword, EnvClientId, EnvClientSecret, EnvEndpoint)
	if err != nil {
		return nil, fmt.Errorf("vimexx: %w", err)
	}

	config := NewDefaultConfig()
	config.baseURL = values[EnvServerBaseURL]
	config.Endpoint = values[EnvEndpoint]
	config.Username = values[EnvUsername]
	config.Password = values[EnvPassword]
	config.ClientId = values[EnvClientId]
	config.ClientSecret = values[EnvClientSecret]

	return NewDNSProviderConfig(config)
}

// NewDNSProviderConfig return a DNSProvider instance configured for Vimexx.
func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("vimexx: the configuration of the DNS provider is nil")
	}

	if config.baseURL == "" || config.Endpoint == "" {
		return nil, errors.New("vimexx: missing server base URL and/or endpoint")
	}

	if config.Username == "" || config.Password == "" || config.ClientId == "" || config.ClientSecret == "" {
		return nil, errors.New("vimexx: incomplete credentials, missing username and/or password and/or client id and/or client secret")
	}

	client := internal.NewClient(config.ClientId, config.ClientSecret, config.Username, config.Password, config.baseURL, config.Endpoint)

	if config.HTTPClient != nil {
		client.HTTPClient = config.HTTPClient
	}

	client.Login()

	return &DNSProvider{
		config: config,
		client: client,
	}, nil
}

// Timeout returns the timeout and interval to use when checking for DNS propagation.
// Adjusting here to cope with spikes in propagation times.
func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

// Present creates a TXT record using the specified parameters.
func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)
	d.client.SetDNSRecord(domain, internal.DNSRecord{Name: info.EffectiveFQDN, Type: "TXT", Content: info.Value, TTL: d.config.TTL})
	return nil
}

// CleanUp removes the TXT record matching the specified parameters.
func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)
	d.client.RemoveDNSRecord(domain, info.EffectiveFQDN, "TXT")
	return nil
}
