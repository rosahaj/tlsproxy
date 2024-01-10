# tlsproxy

Proof-of-concept http proxy for ja3 fingerprint customization and proxy chaining. Not actively mantained.

```
Usage of ./tlsproxy:
  -addr string
        proxy listen address
  -ca string
        tls ca certificate (generated automatically if not present) (default "ca.pem")
  -client-profile string
        default client profile (can be overriden through x-tlsproxy-client-profile header)
  -key string
        tls ca key (generated automatically if not present) (default "key.pem")
  -port string
        proxy listen port (default "8080")
  -upstream-proxy string
        default upstream proxy (can be overriden through x-tlsproxy-upstream-proxy header)
  -verbose
        enable verbose logging
```
