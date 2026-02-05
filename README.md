# WGRest

[![Build Status](https://github.com/suquant/wgrest/actions/workflows/ci.yml/badge.svg)](https://github.com/suquant/wgrest/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/suquant/wgrest/branch/master/graph/badge.svg?token=NM179YJFEJ)](https://codecov.io/gh/suquant/wgrest)

WGRest is a REST API server for WireGuard. It uses wg-quick style configuration files for persistent storage, providing full compatibility with standard WireGuard tooling.

## Features

- **Full wg-quick compatibility** - Uses standard `/etc/wireguard/*.conf` files
- **Device management** - Create, update, delete WireGuard interfaces
- **Peer management** - Full CRUD operations with search and sorting
- **Interface lifecycle** - Bring interfaces up/down via API (`wg-quick up/down`)
- **wg-quick config export** - Download peer configurations as `quick.conf`
- **ACME TLS support** - Automatic Let's Encrypt certificates
- **Bearer token auth** - Simple token-based authorization
- **Swagger UI** - Interactive API documentation at `/docs/`

## Requirements

- Linux with WireGuard kernel module
- `wireguard-tools` package (provides `wg` and `wg-quick`)

## Install

### On Debian / Ubuntu

```shell
curl -L https://github.com/suquant/wgrest/releases/latest/download/wgrest_amd64.deb -o wgrest_amd64.deb
dpkg -i wgrest_amd64.deb
```

### Manual

```shell
curl -L https://github.com/suquant/wgrest/releases/latest/download/wgrest-linux-amd64 -o wgrest
chmod +x wgrest
sudo mv wgrest /usr/local/bin/
```

## Configuration

```shell
wgrest --help

NAME:
   wgrest - REST API for WireGuard

GLOBAL OPTIONS:
   --conf value           wgrest config file path (default: "/etc/wgrest/wgrest.conf")
   --version              Print version and exit
   --listen value         Listen address (default: "127.0.0.1:8000")
   --config-dir value     WireGuard config directory (default: "/etc/wireguard")
   --certs-dir value      ACME TLS certificates cache directory (default: "/var/lib/wgrest/certs")
   --dump-interval value  Config dump interval (default: 10m)
   --static-auth-token value  Bearer token for authorization
   --tls-domain value     TLS Domains for ACME (Let's Encrypt)
   --help, -h             show help
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `WGREST_CONF` | Config file path | `/etc/wgrest/wgrest.conf` |
| `WGREST_LISTEN` | Listen address | `127.0.0.1:8000` |
| `WGREST_CONFIG_DIR` | WireGuard config dir | `/etc/wireguard` |
| `WGREST_CERTS_DIR` | TLS certificates dir | `/var/lib/wgrest/certs` |
| `WGREST_DUMP_INTERVAL` | Config dump interval | `10m` |
| `WGREST_STATIC_AUTH_TOKEN` | Bearer token | - |
| `WGREST_TLS_DOMAIN` | ACME domains | - |

## Quick Start

```shell
# Start server with auth token
wgrest --static-auth-token "secret" --listen "127.0.0.1:8000"

# Open Swagger UI
open http://127.0.0.1:8000/docs/
```

## API Examples

### Create a device

```shell
curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer secret" \
    -d '{
        "name": "wg0",
        "listen_port": 51820,
        "address": ["10.0.0.1/24"]
    }' \
    http://127.0.0.1:8000/v1/devices/
```

### Get devices

```shell
curl -H "Authorization: Bearer secret" \
    http://127.0.0.1:8000/v1/devices/
```

### Update device

```shell
curl -X PATCH \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer secret" \
    -d '{"listen_port": 51821}' \
    http://127.0.0.1:8000/v1/devices/wg0/
```

### Add peer

```shell
curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer secret" \
    -d '{"allowed_ips": ["10.0.0.2/32"]}' \
    http://127.0.0.1:8000/v1/devices/wg0/peers/
```

### Get peers

```shell
curl -H "Authorization: Bearer secret" \
    http://127.0.0.1:8000/v1/devices/wg0/peers/
```

### Bring interface up/down

```shell
# Bring up
curl -X POST \
    -H "Authorization: Bearer secret" \
    http://127.0.0.1:8000/v1/devices/wg0/up/

# Bring down
curl -X POST \
    -H "Authorization: Bearer secret" \
    http://127.0.0.1:8000/v1/devices/wg0/down/
```

### Delete peer

```shell
curl -X DELETE \
    -H "Authorization: Bearer secret" \
    http://127.0.0.1:8000/v1/devices/wg0/peers/{urlSafePubKey}/
```

## URL-Safe Public Keys

Peer public keys in URLs use URL-safe base64 encoding. Convert standard base64:

```python
import base64
pub_key = "hQ1yeyFy+bZn/5jpQNNrZ8MTIGaimZxT6LbWAkvmKjA="
url_safe = base64.urlsafe_b64encode(base64.b64decode(pub_key)).decode()
print(url_safe)  # hQ1yeyFy-bZn_5jpQNNrZ8MTIGaimZxT6LbWAkvmKjA=
```

## Development

```shell
# Build
make build

# Run tests
make test

# Generate swagger docs
make swagger

# Run linter
make lint
```

## Credits

- [ForestVPN.com](https://forestvpn.com) - Free VPN for all
- [SpaceV.net](https://spacev.net) - VPN for teams
