package swarm

import (
	"encoding/json"
	"errors"
	"time"
	"zscaler/core/rule"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// Provider wrapper for docker API client
type Provider struct {
	cli *client.Client
}

// getAPI return a provider
func getAPI() Provider {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}
	return Provider{cli: cli}
}

// Check if service is started
func (sp Provider) Check(rule rule.Default) bool {
	// get the containers tagged with the service
	containers := sp.getTag(rule.Target.Name)
	return len(containers) > 0
}

func (sp Provider) getAll() []types.Container {
	// TODO filtering should be made here
	options := types.ContainerListOptions{All: true}
	// TODO ctx for timeout
	containers, err := sp.cli.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}
	return containers
}

func (sp Provider) getTag(tag string) []types.Container {
	var tagged []types.Container
	containers := sp.getAll()
	for _, c := range containers {
		for _, v := range c.Labels {
			if v == tag {
				tagged = append(tagged, c)
			}
		}
	}
	return tagged
}

func (sp Provider) getStats(cID string) *types.StatsJSON {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := sp.cli.ContainerStats(ctx, cID, false)
	if err != nil {
		return nil
	}
	var stats = new(types.StatsJSON)
	json.NewDecoder(r).Decode(stats)
	return stats
}

// ScaleUp target service
func (sp Provider) ScaleUp(tag string) error {
	containers := sp.getTag(tag)
	if len(containers) == 0 {
		return errors.New("Cannot scale up: target service not found")
	}
	sp.duplicateAndStart(containers[0].ID)
	return nil
}

// duplicateAndStart specifier container
func (sp Provider) duplicateAndStart(cID string) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// r, err = sp.cli.ContainerCreate(ctx, config, hostConfig, networkingConfig, containerName)
	// if err != nil {
	//	return err
	// }
	return nil
}

// ScaleDown target service
func (sp Provider) ScaleDown(tag string) error {
	return nil
}

// Scaler structure
type Scaler struct {
	service string
}

// Describe SwarmScaler
func (s *Scaler) Describe() string {
	return "Swarm scaler for docker"
}

// Up scale
func (s *Scaler) Up() error {
	// launch another instance of the same service
	sp := getAPI()
	sp.ScaleUp(s.service)
	return nil
}

// Down scale
func (s *Scaler) Down() error {
	sp := getAPI()
	sp.ScaleDown(s.service)
	return nil
}
