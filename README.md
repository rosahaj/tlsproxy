# tlsproxy

HTTP proxy with per-request uTLS fingerprint mimicry and upstream proxy tunneling. Currently WIP.

Built on top of [uTLS](https://github.com/refraction-networking/utls) and [goproxy](https://github.com/elazarl/goproxy/). Inspired by [ja3proxy](https://github.com/LyleMi/ja3proxy).

## Usage

### Building from source

```bash
git clone https://github.com/rosahaj/tlsproxy
cd ja3proxy
go build

# Start proxy
./tlsproxy -client Chrome-120
# Make requests
curl --cacert cert.pem --proxy http://localhost:8080 https://www.example.com
```

Pre-built binaries are available in the [Releases](https://github.com/rosahaj/tlsproxy/releases) section.

### Using docker CLI

```bash
docker run \
      -v ./credentials:/app/credentials \
      -p 8080:8080 \
      ghcr.io/rosahaj/tlsproxy:latest \
      -cert /app/credentials/cert.pem \
      -key /app/credentials/key.pem \
      -client Chrome-120
```

### Using docker compose

See [`compose.yaml`](https://github.com/rosahaj/tlsproxy/blob/master/compose.yaml)

```bash
docker compose up -d
```

### CLI usage

```
Usage of ./tlsproxy:
  -addr string
        Proxy listen address
  -cert string
        TLS CA certificate (generated automatically if not present) (default "cert.pem")
  -client string
        Default utls clientHelloID (can be overriden through x-tlsproxy-client header) (default "Chrome-120")
  -key string
        TLS CA key (generated automatically if not present) (default "key.pem")
  -port string
        Proxy listen port (default "8080")
  -upstream string
        Default upstream proxy (can be overriden through x-tlsproxy-upstream header)
  -verbose
        Enable verbose logging
```