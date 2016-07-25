package probe

import (
	"fmt"
	"testing"
	"time"
)

func TestDefaultScalingProbe(t *testing.T) {
	p := new(DefaultScalingProbe)
	for i := 1; i <= 10; i++ {
		val, _ := p.Value()
		fmt.Printf("%f\n", val)
		time.Sleep(2 * time.Second)
	}
}
