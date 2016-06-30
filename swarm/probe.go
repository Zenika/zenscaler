package swarm

import "github.com/docker/engine-api/types"

// AverageCPU probe of all containers
type AverageCPU struct {
	Tag string
}

// Name of probe
func (info AverageCPU) Name() string {
	return "AverageCPU"
}

// Value return average CPU consuption of all service containers
func (info AverageCPU) Value() float64 {
	sp := getAPI()
	containers := sp.getTag(info.Tag)

	var cpusum float64
	for _, c := range containers {
		cpusum = +calculateCPUPercent(sp.getStats(c.ID))
	}
	return cpusum / float64(len(containers))
}

func calculateCPUPercent(v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}
