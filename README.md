# wgrest
WireGuard REST API

WireGuard is an simple and modern VPN. It is cross-platform (Windows, macOS, BSD, iOS, Android). 


## Install

```shell
curl -L https://github.com/suquant/wgrest/releases/download/v0.0.1/wgrest-linux-amd64 -o wgrest

chmod +x wgrest
```

## Run WireGuard REST API Server

```shell
wgrest --token=secret --scheme=http --port=8000
```

## Create **wg0** device

```shell
curl -v -g \
    -H "Accept: */*" \
    -H "Content-Type: application/json" \
    -H "Token: secret" \
    -X POST \
    -d '{
        "name": "wg0", 
        "listen_port":51820, 
        "private_key": "cLmxIyJx/PGWrQlevBGr2LQNOqmBGYbVfu4XcRO2SEo=", 
        "network": "10.10.1.1/24"
    }' \
    http://127.0.0.1:8000/devices/
```

## Get devices

```shell
curl -v -g \
    -H "Accept: */*" \
    -H "Content-Type: application/json" \
    -H "Token: secret" \
    -X GET \
    http://127.0.0.1:8000/devices/
```

## Add peer

```shell
curl -v -g \
    -H "Accept: */*" \
    -H "Content-Type: application/json" \
    -H "Token: secret" \
    -X POST \
    -d '{
        "public_key": "hQ1yeyFy+bZn/5jpQNNrZ8MTIGaimZxT6LbWAkvmKjA=", 
        "allowed_ips": ["10.10.1.2/32"], 
        "preshared_key": "uhFI9c9rInyxqgZfeejte6apHWbewoiy32+Bo34xRFs="
    }' \
    http://127.0.0.1:8000/devices/wg0/peers/
```

## Get peers

```shell
curl -v -g \
    -H "Accept: */*" \
    -H "Content-Type: application/json" \
    -H "Token: secret" \
    -X GET \
    http://127.0.0.1:8000/devices/wg0/peers/
```

## Delete peer

Since the wireguard public key is the standard base64 encoded string, it is not safe to use in URI schema,
is that reason peer_id contains the same public key of the peer but encoded with URL safe base64 encoder.

peer_id can be retrieved either by `peer_id` field from peer list endpoint or by this rule

```shell
python3 -c "import base64; \
    print(\
        base64.urlsafe_b64encode(\
            base64.b64decode('hQ1yeyFy+bZn/5jpQNNrZ8MTIGaimZxT6LbWAkvmKjA=')\
        ).decode()\
    )"
```

delete peer request

```shell
curl -v -g \
    -H "Accept: */*" \
    -H "Content-Type: application/json" \
    -H "Token: secret" \
    -X DELETE \
    http://127.0.0.1:8000/devices/wg0/peers/hQ1yeyFy-bZn_5jpQNNrZ8MTIGaimZxT6LbWAkvmKjA=/
```