// Package core provide interface definition and execution flow
package core

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Zenika/zscaler/core/types"
)

// Config store the current running configuration
var Config *types.Configuration

// CheckTLS for missing certificate and key
func CheckTLS() (bool, error) {
	o := Config.Orchestrator
	if o.TLSCACertPath+o.TLSCertPath+o.TLSKeyPath == "" {
		return false, nil
	}
	switch "" {
	case o.TLSCACertPath:
		return false, fmt.Errorf("tls-cacert path not provided")
	case o.TLSCertPath:
		return false, fmt.Errorf("tls-cert path not provided")
	case o.TLSKeyPath:
		return false, fmt.Errorf("tls-key path not provided")
	}
	// all set, TLS seems ok !
	return true, nil
}

// HTTPSClient return configured http client
func HTTPSClient() (*http.Client, error) {
	o := Config.Orchestrator
	//try to load up files...
	ca, err := ioutil.ReadFile(o.TLSCACertPath)
	if err != nil {
		return nil, fmt.Errorf("tls-cacert: %s", err)
	}
	cert, err := ioutil.ReadFile(o.TLSCertPath)
	if err != nil {
		return nil, fmt.Errorf("tls-cert: %s", err)
	}
	key, err := ioutil.ReadFile(o.TLSKeyPath)
	if err != nil {
		return nil, fmt.Errorf("tls-key: %s", err)
	}
	// load PEM
	certPair, err := tls.LoadX509KeyPair(string(cert), string(key))
	if err != nil {
		return nil, fmt.Errorf("cannot load key pair: %s", err)
	}
	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("failed to load CA file")
	}
	// build up https configuration
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: append(make([]tls.Certificate, 1, 1), certPair),
			RootCAs:      caPool,
		},
	}
	return &http.Client{Transport: tr}, nil
}
