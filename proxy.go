package main

import (
	"flag"
	"log"
	"net/http"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/elazarl/goproxy"
)

func main() {
	flag.StringVar(&Flags.addr, "addr", "", "proxy listen address")
	flag.StringVar(&Flags.port, "port", "8080", "proxy listen port")
	flag.StringVar(&Flags.ca, "ca", "ca.pem", "tls ca certificate (generated automatically if not present)")
	flag.StringVar(&Flags.key, "key", "key.pem", "tls ca key (generated automatically if not present)")
	flag.StringVar(&Flags.upstreamProxy, "upstream-proxy", "", "default upstream proxy (can be overriden through x-tlsproxy-upstream-proxy header)")
	flag.StringVar(&Flags.clientProfile, "client-profile", "", "default client profile (can be overriden through x-tlsproxy-client-profile header)")
	flag.BoolVar(&Flags.verbose, "verbose", false, "enable verbose logging")
	flag.Parse()

	loadDefaultProxyConfig()
	loadCA()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = Flags.verbose

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			proxyConfig := parseCustomHeaders(&req.Header)
			removeCustomHeaders(&req.Header)

			ctx.RoundTripper = goproxy.RoundTripperFunc(
				func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Response, error) {
					options := []tls_client.HttpClientOption{
						tls_client.WithTimeoutSeconds(30),
						tls_client.WithNotFollowRedirects(),
					}

					if len(proxyConfig.upstreamProxy) > 0 {
						if isValidUrl(proxyConfig.upstreamProxy) {
							options = append(options, tls_client.WithProxyUrl(proxyConfig.upstreamProxy))
						} else {
							return invalidUpstreamProxyResponse(req, ctx, proxyConfig.upstreamProxy), nil
						}
					} else if len(DefaultUpstreamProxy) > 0 {
						options = append(options, tls_client.WithProxyUrl(DefaultUpstreamProxy))
					}

					if len(proxyConfig.clientProfile) > 0 {
						clientProfile, ok := profiles.MappedTLSClients[proxyConfig.clientProfile]
						if !ok {
							return invalidClientProfileResponse(req, ctx, proxyConfig.clientProfile), nil
						}

						options = append(options, tls_client.WithClientProfile(clientProfile))
					} else if DefaultClientProfile.isSet {
						options = append(options, tls_client.WithClientProfile(DefaultClientProfile.clientProfile))
					}

					var logger tls_client.Logger
					if Flags.verbose {
						logger = tls_client.NewDebugLogger(tls_client.NewLogger())
					} else {
						logger = tls_client.NewNoopLogger()
					}

					client, err := tls_client.NewHttpClient(logger, options...)
					if err != nil {
						log.Fatal("Unable to create tls_client client", err)
					}

					request, err := convertHttpReqToFhttpReq(req)
					if err != nil {
						log.Fatal("Unable to convert request to fhttp Request", err)
					}

					response, err := client.Do(request)

					resp := convertFhttpRespToHttpResp(response, req)

					return resp, err
				},
			)
			return req, nil
		},
	)

	listenAddr := Flags.addr + ":" + Flags.port
	log.Println("tlsproxy listening at " + listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, proxy))

	// time.Sleep(1 * time.Second)

	// proxyUrl := "http://localhost:8081"
	// request, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
	// request.Header.Add("x-tlsproxy-client-profile", "chrome_117")
	// if err != nil {
	// 	log.Fatalf("new request failed:%v", err)
	// }
	// certPool := x509.NewCertPool()
	// var x509ca *x509.Certificate
	// if x509ca, err = x509.ParseCertificate(goproxy.GoproxyCa.Certificate[0]); err != nil {
	// 	return
	// }
	// certPool.AddCert(x509ca)
	// tr := &http.Transport{Proxy: func(req *http.Request) (*url.URL, error) { return url.Parse(proxyUrl) }, TLSClientConfig: &tls.Config{RootCAs: certPool}}
	// client := &http.Client{Transport: tr}
	// rsp, err := client.Do(request)
	// if err != nil {
	// 	log.Fatalf("get rsp failed:%v", err)

	// }
	// defer rsp.Body.Close()
	// data, _ := ioutil.ReadAll(rsp.Body)

	// if rsp.StatusCode != http.StatusOK {
	// 	log.Fatalf("status %d, data %s", rsp.StatusCode, data)
	// }

	// log.Printf("norm_id: %s", data)
}
