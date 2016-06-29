package probe

import "time"

// Probe interface
type Probe interface {
	Name() string
	Value() float32
}

// DefaultScalingProbe report a fake sensor value
// Value goes from 0 to 1 and to 1 to 0 each minute
type DefaultScalingProbe struct{}

// Name of the probe
func (p *DefaultScalingProbe) Name() string {
	return "DefaultScalingProbe"
}

// Value of the probe
func (p *DefaultScalingProbe) Value() float32 {
	_, _, s := time.Now().Clock()
	return float32(abs(s-30)) / float32(30)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// Initialize some defaul probes
func Initialize() map[string]Probe {
	return map[string]Probe{
		"DefaultScalingProbe": new(DefaultScalingProbe),
	}
}
