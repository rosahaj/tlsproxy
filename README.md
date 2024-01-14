# tlsproxy

HTTP proxy with per-request uTLS fingerprint mimicry and upstream proxy tunneling. Currently WIP.

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