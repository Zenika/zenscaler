package probe

import (
	"fmt"
	"testing"
	"time"
)

func TestDefaultScalingProbe(t *testing.T) {
	p := new(DefaultScalingProbe)
	for i := 1; i <= 10; i++ {
		fmt.Printf("%f\n", p.Value())
		time.Sleep(2 * time.Second)
	}
}
