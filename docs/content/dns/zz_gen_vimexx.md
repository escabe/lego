---
title: "vimexx.nl"
date: 2019-03-03T16:39:46+01:00
draft: false
slug: vimexx
dnsprovider:
  since:    "v4.11.0"
  code:     "vimexx"
  url:      "https://www.vimexx.nl/"
---

<!-- THIS DOCUMENTATION IS AUTO-GENERATED. PLEASE DO NOT EDIT. -->
<!-- providers/dns/vimexx/vimexx.toml -->
<!-- THIS DOCUMENTATION IS AUTO-GENERATED. PLEASE DO NOT EDIT. -->


Configuration for [vimexx.nl](https://www.vimexx.nl/).


<!--more-->

- Code: `vimexx`
- Since: v4.11.0


Here is an example bash command using the vimexx.nl provider:

```bash
VIMEXX_SERVER_BASE_URL="https://api.vimexx.nl" \
VIMEXX_ENDPOINT="/api/v1" \
VIMEXX_USERNAME="you@example.com" \
VIMEXX_PASSWORD=supersecret \
VIMEXX_CLIENT_ID="123" \
VIMEXX_CLIENT_SECRET="abcdef" \
lego --email you@example.com --dns vimexx -d '*.example.com' run
```




## Credentials

| Environment Variable Name | Description |
|-----------------------|-------------|
| `VIMEXX_CLIENT_ID` | Vimexx client id |
| `VIMEXX_CLIENT_SECRET` | Vimexx client secret |
| `VIMEXX_ENDPOINT` | Base API (ex: /api/v1 or /apitest/v1) |
| `VIMEXX_PASSWORD` | Vimexx password |
| `VIMEXX_SERVER_BASE_URL` | Base URL of the server (ex: https://api.vimexx.nl) |
| `VIMEXX_USERNAME` | Vimexx username (e-mail) |

The environment variable names can be suffixed by `_FILE` to reference a file instead of a value.
More information [here]({{% ref "dns#configuration-and-credentials" %}}).


## Additional Configuration

| Environment Variable Name | Description |
|--------------------------------|-------------|
| `VIMEXX_HTTP_TIMEOUT` | API request timeout in seconds (Default: 30) |
| `VIMEXX_POLLING_INTERVAL` | Time between DNS propagation check in seconds (Default: 2) |
| `VIMEXX_PROPAGATION_TIMEOUT` | Maximum waiting time for DNS propagation in seconds (Default: 60) |
| `VIMEXX_TTL` | The TTL of the TXT record used for the DNS challenge in seconds (Default: 300) |

The environment variable names can be suffixed by `_FILE` to reference a file instead of a value.
More information [here]({{% ref "dns#configuration-and-credentials" %}}).




## More information



<!-- THIS DOCUMENTATION IS AUTO-GENERATED. PLEASE DO NOT EDIT. -->
<!-- providers/dns/vimexx/vimexx.toml -->
<!-- THIS DOCUMENTATION IS AUTO-GENERATED. PLEASE DO NOT EDIT. -->
