package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	cflog "github.com/cloudflare/cfssl/log"
	"github.com/elazarl/goproxy"
	utls "github.com/refraction-networking/utls"
	sf "gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/snowflake/v2/common/utls"
)

func main() {
	flag.StringVar(&Flags.addr, "addr", "", "Proxy listen address")
	flag.StringVar(&Flags.port, "port", "8080", "Proxy listen port")
	flag.StringVar(&Flags.cert, "cert", "cert.pem", "TLS CA certificate (generated automatically if not present)")
	flag.StringVar(&Flags.key, "key", "key.pem", "TLS CA key (generated automatically if not present)")
	flag.StringVar(&Flags.upstreamProxy, "upstream", "", "Default upstream proxy prefixed by \"socks5://\" (can be overriden through x-tlsproxy-upstream header)")
	flag.StringVar(&Flags.client, "client", "Chrome-120", "Default utls clientHelloID (can be overriden through x-tlsproxy-client header)")
	flag.BoolVar(&Flags.verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	if !Flags.verbose {
		cflog.Level = cflog.LevelError
	}

	loadDefaultProxyConfig()
	loadCA()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = Flags.verbose

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			proxyConfig := parseCustomHeaders(&req.Header)
			removeCustomHeaders(&req.Header)

			clientHelloId := DefaultClientHelloID
			upstreamProxy := DefaultUpstreamProxy

			if len(proxyConfig.client) > 0 {
				customClientHeaderId, ok := getClientHelloID(proxyConfig.client)
				if !ok {
					return req, invalidClientResponse(req, ctx, proxyConfig.client)
				}

				clientHelloId = customClientHeaderId
			}

			if len(proxyConfig.upstreamProxy) > 0 {
				proxyUrl, err := url.Parse(proxyConfig.upstreamProxy)
				if err != nil {
					return req, invalidUpstreamProxyResponse(req, ctx, proxyConfig.upstreamProxy)
				}

				upstreamProxy = proxyUrl
			}

			roundTripper := sf.NewUTLSHTTPRoundTripperWithProxy(clientHelloId, &utls.Config{
				InsecureSkipVerify: true,
			}, http.DefaultTransport, false, upstreamProxy)

			ctx.RoundTripper = goproxy.RoundTripperFunc(
				func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Response, error) {
					return roundTripper.RoundTrip(req)
				})

			return req, nil
		},
	)

	listenAddr := Flags.addr + ":" + Flags.port
	log.Println("tlsproxy listening at " + listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, proxy))
}
