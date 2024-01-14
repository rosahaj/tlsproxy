package main

import "net/http"

const (
	UpStreamProxyHeader = "x-tlsproxy-upstream"
	ClientProfileHeader = "x-tlsproxy-client"
)

var CustomHeaders = []string{UpStreamProxyHeader, ClientProfileHeader}

type ProxyConfig struct {
	client        string
	upstreamProxy string
}

func parseCustomHeaders(headers *http.Header) ProxyConfig {
	return ProxyConfig{
		upstreamProxy: headers.Get(UpStreamProxyHeader),
		client:        headers.Get(ClientProfileHeader),
	}
}

func removeCustomHeaders(headers *http.Header) {
	for _, header := range CustomHeaders {
		headers.Del(header)
	}
}
