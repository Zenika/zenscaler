package provider

import (
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

func Init() *client.Client {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}
	return cli
}

func getTaggedContainers(cli *client.Client) []types.Container {
	options := types.ContainerListOptions{All: true}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}
	return containers
}
