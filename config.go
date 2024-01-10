package main

import (
	"log"
	"net/url"

	"github.com/bogdanfinn/tls-client/profiles"
)

type CLIFlags struct {
	addr          string
	port          string
	ca            string
	key           string
	upstreamProxy string
	clientProfile string
	verbose       bool
}

var (
	Flags                CLIFlags
	DefaultClientProfile = &struct {
		isSet         bool
		clientProfile profiles.ClientProfile
	}{}
	DefaultUpstreamProxy string
)

func setDefaultClientProfile(clientProfile string) {
	if loadedClientProfile, ok := profiles.MappedTLSClients[clientProfile]; ok {
		DefaultClientProfile.isSet = true
		DefaultClientProfile.clientProfile = loadedClientProfile
	} else {
		log.Fatalf("Invalid client profile: %s", clientProfile)
	}
}

func isValidUrl(inputUrl string) bool {
	_, err := url.Parse(inputUrl)
	return err == nil
}

func setDefaultUpstreamProxy(upstreamProxy string) {
	if isValidUrl(upstreamProxy) {
		DefaultUpstreamProxy = upstreamProxy
	} else {
		log.Fatalf("Invalid upstream proxy: %s", upstreamProxy)
	}
}

func loadDefaultProxyConfig() {
	if len(Flags.clientProfile) > 0 {
		setDefaultClientProfile(Flags.clientProfile)
	}

	if len(Flags.upstreamProxy) > 0 {
		setDefaultUpstreamProxy(Flags.upstreamProxy)
	}
}
