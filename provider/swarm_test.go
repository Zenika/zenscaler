package provider

import (
	"fmt"
	"testing"
)

func TestSwarmProvider(t *testing.T) {
	cli := Init()
	containers := getTaggedContainers(cli)

	for _, c := range containers {
		fmt.Println("**" + c.Names[0] + "**")
		for k, v := range c.Labels {
			fmt.Println(k + ":" + v)
		}
	}
}
