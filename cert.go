package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/elazarl/goproxy"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func setGoproxyCA(tlsCert tls.Certificate) {
	var err error
	if tlsCert.Leaf, err = x509.ParseCertificate(tlsCert.Certificate[0]); err != nil {
		log.Fatal("Unable to parse ca", err)
	}

	goproxy.GoproxyCa = tlsCert
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
}

func loadCA() {
	if fileExists(Flags.ca) && fileExists(Flags.key) {
		tlsCert, err := tls.LoadX509KeyPair(Flags.ca, Flags.key)
		if err != nil {
			log.Fatal("Unable to load ca", err)
		}

		setGoproxyCA(tlsCert)
	} else {
		if fileExists(Flags.ca) {
			log.Fatalf("CA certificate exists, but found no corresponding key at %s", Flags.key)
		} else if fileExists(Flags.key) {
			log.Fatalf("CA key exists, but found no corresponding certificate at %s", Flags.ca)
		}

		// TODO
		// log.Println("CA does not exist, generating certificate and key")
		log.Println("No CA certificate and key found, please generate")
		os.Exit(-1)
	}
}
