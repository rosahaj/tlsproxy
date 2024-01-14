package main

import (
	"log"
	"net/url"
	"strings"

	utls "github.com/refraction-networking/utls"
)

type CLIFlags struct {
	addr          string
	port          string
	cert          string
	key           string
	upstreamProxy string
	client        string
	verbose       bool
}

var (
	Flags                CLIFlags
	DefaultClientHelloID utls.ClientHelloID
	DefaultUpstreamProxy *url.URL
)

func getClientHelloID(client string) (utls.ClientHelloID, bool) {
	clientArr := strings.Split(client, "-")
	if len(clientArr) != 2 {
		return utls.ClientHelloID{}, false
	}

	return utls.ClientHelloID{
		Client:  clientArr[0],
		Version: clientArr[1],
		Seed:    nil,
		Weights: nil,
	}, true
}

func setDefaultClientHelloID(client string) {
	clientHelloId, ok := getClientHelloID(client)
	if !ok {
		log.Fatalf("Invalid client format: %s", client)
	}

	DefaultClientHelloID = clientHelloId
}

func setDefaultUpstreamProxy(upstreamProxy string) {
	proxyUrl, err := url.Parse(upstreamProxy)
	if err != nil {
		log.Fatalf("Invalid upstream proxy: %s", upstreamProxy)
	}

	DefaultUpstreamProxy = proxyUrl
}

func loadDefaultProxyConfig() {
	setDefaultClientHelloID(Flags.client)

	if len(Flags.upstreamProxy) > 0 {
		setDefaultUpstreamProxy(Flags.upstreamProxy)
	}
}
