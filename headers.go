package main

import "net/http"

const (
	UpStreamProxyHeader = "x-tlsproxy-upstream-proxy"
	ClientProfileHeader = "x-tlsproxy-client-profile"
)

var CustomHeaders = []string{UpStreamProxyHeader, ClientProfileHeader}

type ProxyConfig struct {
	upstreamProxy string
	clientProfile string
}

func parseCustomHeaders(headers *http.Header) ProxyConfig {
	return ProxyConfig{
		upstreamProxy: headers.Get(UpStreamProxyHeader),
		clientProfile: headers.Get(ClientProfileHeader),
	}
}

func removeCustomHeaders(headers *http.Header) {
	for _, header := range CustomHeaders {
		headers.Del(header)
	}
}
