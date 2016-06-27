package probe

import (
	"time"
)

// Probe interface
type Probe interface {
	GetName() string
	GetValue() float32
}

// DefaultScalingProbe report a fake sensor value
// Value goes from 0 to 1 and to 1 to 0 each minute
type DefaultScalingProbe struct{}

// GetName of the scaler
func (p *DefaultScalingProbe) GetName() string {
	return "DefaultScalingProbe"
}

// GetValue of the scaler
func (p *DefaultScalingProbe) GetValue() float32 {
	_, _, s := time.Now().Clock()
	return float32(abs(s-30)) / float32(30)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
