package scaler

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/Zenika/zscaler/core"
	"github.com/Zenika/zscaler/core/tls"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// ServiceScaler work with docker 1.12 swarm services (API 1.24)
type ServiceScaler struct {
	ServiceID    string `json:"service"`
	EngineSocket string `json:"socket"`
	cli          *client.Client
}

// Describe scaler
func (s *ServiceScaler) Describe() string {
	return "Docker 1.12 swarm mode API scaler"
}

// JSON encode
func (s *ServiceScaler) JSON() ([]byte, error) {
	encoded, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// Up using API on swarm socket
func (s *ServiceScaler) Up() error {
	err := s.scaleService(func(n uint64) uint64 {
		return n + 1
	})
	return err
}

// Down using API on swarm socket
func (s *ServiceScaler) Down() error {
	err := s.scaleService(func(n uint64) uint64 {
		if n > 1 {
			return n - 1
		}
		return n
	})
	return err
}

// Update service target replicas
func (s *ServiceScaler) scaleService(scale func(uint64) uint64) error {
	cli, err := s.getDocker()
	if err != nil {
		return err
	}
	ctx := context.Background()
	service, _, err := cli.ServiceInspectWithRaw(ctx, s.ServiceID)
	if err != nil {
		return err
	}
	serviceMode := &service.Spec.Mode
	if serviceMode.Replicated == nil {
		return fmt.Errorf("scale can only be used with replicated mode")
	}
	target := scale(*serviceMode.Replicated.Replicas)
	log.WithFields(log.Fields{
		"service": s.ServiceID,
		"count":   *serviceMode.Replicated.Replicas,
		"target":  target,
	}).Debugf("scale service")
	serviceMode.Replicated.Replicas = &target

	err = cli.ServiceUpdate(ctx, service.ID, service.Version, service.Spec, types.ServiceUpdateOptions{})
	return err
}

// Lazy init of new API channel to docker engine
func (s *ServiceScaler) getDocker() (cli *client.Client, err error) {
	var HTTPClient *http.Client
	if core.Config.Orchestrator.TLS {
		HTTPClient, err = tls.HTTPSClient()
		if err != nil {
			return nil, err
		}
	}
	if s.cli == nil {
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		s.cli, err = client.NewClient(s.EngineSocket, "v1.24", HTTPClient, defaultHeaders)
	}
	return s.cli, err
}
