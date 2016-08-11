// Package core provide interface definition and execution flow
package core

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Zenika/zscaler/core/rule"
	"github.com/Zenika/zscaler/core/scaler"

	log "github.com/Sirupsen/logrus"
)

const bufferSize = 10

// Config store the current running configuration
var Config *Configuration

// Configuration holder
type Configuration struct {
	Orchestrator OrchestratorConfig
	Scalers      map[string]scaler.Scaler
	Rules        map[string]rule.Rule
	Errchan      chan error
}

// OrchestratorConfig hold all necessary connection informations
type OrchestratorConfig struct {
	Kind          string // docker, kubernetes, mesos...
	Endpoint      string // http adress:port or unix socket
	TLSCACertPath string // ca cert path, PEM formated
	TLSCertPath   string // user cert path, PEM formated
	TLSKeyPath    string // user private key path, PEM formated
}

// CheckTLS for missing certificate and key
func (o OrchestratorConfig) CheckTLS() error {
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
	return nil
}

// HTTPSClient return configured http client
func (o OrchestratorConfig) HTTPSClient() (*http.Client, error) {
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

// Initialize core module
func (c Configuration) Initialize() {
	c.Errchan = make(chan error, bufferSize)
	c.loop()
}

// event loop
func (c Configuration) loop() {
	log.Debug("Enter control loop...")
	// lanch a watcher on each rule
	for _, r := range c.Rules {
		go rule.Watcher(c.Errchan, r)
	}
	// watch for errors
	for {
		err := <-c.Errchan
		if err != nil {
			log.Errorf("%s", err)
		}
	}
}
