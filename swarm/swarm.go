package swarm

import (
	"encoding/json"
	"time"
	"zscaler/core/rule"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

// Provider wrapper for docker API client
type Provider struct {
	cli *client.Client
}

var provider *Provider

// getAPI return a lazy-initialized provider
func getAPI() Provider {
	if provider == nil {
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		cli, err := client.NewClient(viper.GetString("endpoint"), "v1.22", nil, defaultHeaders)
		if err != nil {
			log.Panic(err)
		}
		provider = &Provider{cli: cli}
	}
	return *provider
}

// Check if service is started
func (sp Provider) Check(rule rule.Default) bool {
	// get the containers tagged with the service
	containers := sp.getTag(rule.ServiceName)
	return len(containers) > 0
}

func (sp Provider) getAll() []types.Container {
	// TODO filtering should be made here
	options := types.ContainerListOptions{All: true}
	// TODO ctx for timeout
	containers, err := sp.cli.ContainerList(context.Background(), options)
	if err != nil {
		log.Panic(err)
	}
	return containers
}

func (sp Provider) getTag(tag string) []types.Container {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tagFilter := filters.NewArgs()
	tagFilter.Add("label", "com.docker.compose.service="+tag)
	options := types.ContainerListOptions{Filter: tagFilter}
	tagged, err := sp.cli.ContainerList(ctx, options)
	if err != nil {
		log.Panic(err)
	}
	return tagged
}

func (sp Provider) getStats(cID string) *types.StatsJSON {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Debug("Get Stats for " + cID)
	r, err := sp.cli.ContainerStats(ctx, cID, false)
	if err != nil {
		log.Errorf("%s", err)
	}
	var stats = new(types.StatsJSON)
	err = json.NewDecoder(r).Decode(stats)
	if err != nil {
		log.Errorf("%s", err)
	}
	err = r.Close()
	if err != nil {
		log.Errorf("%s", err)
	}
	return stats
}
