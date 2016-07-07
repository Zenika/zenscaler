package swarm

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"zscaler/core/rule"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"
)

// Provider wrapper for docker API client
type Provider struct {
	cli *client.Client
}

var provider *Provider

// getAPI return a provider
func getAPI() Provider {
	if provider == nil {
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
		if err != nil {
			panic(err)
		}
		provider = &Provider{cli: cli}
	}
	return *provider
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tagFilter := filters.NewArgs()
	tagFilter.Add("label", "com.docker.compose.service="+tag)
	options := types.ContainerListOptions{Filter: tagFilter}
	fmt.Println("Get ContainerList")
	tagged, err := sp.cli.ContainerList(ctx, options)
	if err != nil {
		panic(err)
	}
	fmt.Println("Got it !")
	return tagged
}

func (sp Provider) getStats(cID string) *types.StatsJSON {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("Get Stats for " + cID)
	r, err := sp.cli.ContainerStats(ctx, cID, false)
	if err != nil {
		return nil
	}
	fmt.Println("Decode JSON for " + cID)
	var stats = new(types.StatsJSON)
	json.NewDecoder(r).Decode(stats)
	r.Close()
	fmt.Println("return stats for " + cID)
	return stats
}

// ScaleUp target service
func (sp Provider) ScaleUp(tag string) error {
	containers := sp.getTag(tag)
	if len(containers) == 0 {
		return errors.New("Cannot scale up: target [" + tag + "] not found")
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
