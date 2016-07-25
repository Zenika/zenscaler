package probe

import "time"

// Probe interface
type Probe interface {
	Name() string
	Value() (float64, error)
}

// DefaultScalingProbe report a fake sensor value
// Value goes from 0 to 1 and to 1 to 0 each minute
type DefaultScalingProbe struct{}

// Name of the probe
func (p *DefaultScalingProbe) Name() string {
	return "DefaultScalingProbe"
}

// Value of the probe
func (p *DefaultScalingProbe) Value() (float64, error) {
	_, _, s := time.Now().Clock()
	return float64(abs(s-30)) / float64(30), nil
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
