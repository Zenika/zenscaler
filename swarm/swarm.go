package swarm

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"

	"github.com/Zenika/zenscaler/core"
	"github.com/Zenika/zenscaler/core/tls"
)

// Provider wrapper for docker API client
type Provider struct {
	cli *client.Client
}

var provider *Provider

// getAPI return a lazy-initialized provider
func getAPI() Provider {
	if provider == nil {
		var HTTPClient *http.Client
		if core.Config.Orchestrator.TLS {
			var err error
			HTTPClient, err = tls.HTTPSClient()
			if err != nil {
				log.Panic(err)
			}
		}
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		cli, err := client.NewClient(core.Config.Orchestrator.Endpoint, "v1.22", HTTPClient, defaultHeaders)
		if err != nil {
			log.Panic(err)
		}
		provider = &Provider{cli: cli}
	}
	return *provider
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
	log.WithFields(log.Fields{
		"cid": cID[0:12],
		"tls": core.Config.Orchestrator.TLS,
	}).Debugf("swarm: get stats")
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
