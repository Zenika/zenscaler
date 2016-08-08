package probe

import (
	"fmt"
	"testing"
)

func TestHAproxyRTimeProbe(t *testing.T) {
	haprobe := &HAproxy{
		Socket: "haproxy.stats",
		Type:   "backend",
		Item:   "rtime",
	}
	val, err := haprobe.Value()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("rtime: %fms\n", val)
}
