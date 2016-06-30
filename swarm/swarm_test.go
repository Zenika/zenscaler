package swarm

import (
	"fmt"
	"testing"
)

func TestDockerGetTags(t *testing.T) {
	containers := getAPI().getTag("traefik")

	for _, c := range containers {
		for _, n := range c.Names {
			fmt.Println(n)
		}
		for k, v := range c.Labels {
			fmt.Println("- [" + k + "] " + v)
		}
	}
}

func TestPobe(t *testing.T) {
	var probe = AverageCPU{
		tag: "traefik",
	}
	fmt.Printf("%.2f\n", probe.Value())
}
