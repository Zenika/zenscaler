package swarm

import (
	"zscaler/core/service"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

type SwarmProvider struct {
	cli *client.Client
}

func getAPI() *client.Client {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}
	return cli
}

// check if services are started on the docker host and get their current state
func CheckServices(services []service.Service) []service.Service {
	ok := make([]service.Service, 0)
	// get the containers tagged with the service
	for _, s := range services {
		if len(getTag(s.Name)) > 0 {
			ok = append(ok, s)
		}
	}
	return ok
}

func getAll() []types.Container {
	cli := getAPI()
	// TODO filtering should be made here
	options := types.ContainerListOptions{All: true}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}
	return containers
}

func getTag(tag string) []types.Container {
	tagged := make([]types.Container, 0)
	containers := getAll()
	for _, c := range containers {
		for _, v := range c.Labels {
			if v == tag {
				tagged = append(tagged, c)
			}
		}
	}
	return tagged
}
