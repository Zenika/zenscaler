package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/Zenika/zenscaler/core"
	"github.com/Zenika/zenscaler/core/types"
)

// CheckTLS for missing certificate and key
func CheckTLS(config *types.Configuration) error {
	o := &config.Orchestrator
	o.TLS = false // enforcing default
	if o.TLSCACertPath+o.TLSCertPath+o.TLSKeyPath == "" {
		return nil
	}
	switch "" {
	case o.TLSCACertPath:
		return fmt.Errorf("tls-cacert path not provided")
	case o.TLSCertPath:
		return fmt.Errorf("tls-cert path not provided")
	case o.TLSKeyPath:
		return fmt.Errorf("tls-key path not provided")
	}
	// all set, TLS seems ok !
	o.TLS = true
	return nil
}

// HTTPSClient return a http client with TLS
func HTTPSClient() (*http.Client, error) {
	o := core.Config.Orchestrator
	//try to load up files...
	ca, err := ioutil.ReadFile(o.TLSCACertPath)
	if err != nil {
		return nil, fmt.Errorf("tls-cacert: %s", err)
	}
	// load PEM
	certPair, err := tls.LoadX509KeyPair(o.TLSCertPath, o.TLSKeyPath)
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

// CheckTLSConfigPath for docker-compose cli compatibility
//
// As docker-compose don't allow to specify each file path, they need to be
// in the same folder and be named ca.pem, cert.pem and key.pem. They must bear
// no password protection.
func CheckTLSConfigPath() (certsPath string, err error) {
	caPath, ca := path.Split(core.Config.Orchestrator.TLSCACertPath)
	if ca != "ca.pem" {
		return "", fmt.Errorf("ca file must be named ca.pem")
	}
	certPath, cert := path.Split(core.Config.Orchestrator.TLSCertPath)
	if cert != "cert.pem" {
		return "", fmt.Errorf("cert file must be named cert.pem")
	}
	keyPath, key := path.Split(core.Config.Orchestrator.TLSKeyPath)
	if key != "key.pem" {
		return "", fmt.Errorf("key file must be named key.pem")
	}
	if caPath != certPath || certPath != keyPath {
		return "", fmt.Errorf("ca.pem, key.pem and cert.pem must be in the same folder")
	}
	return caPath, nil
}
