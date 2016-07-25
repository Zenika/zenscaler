package probe

import (
	"fmt"
	"testing"
)

func TestHAproxyRTimeProbe(t *testing.T) {
	haprobe := &HAproxy{
		Socket: "/home/maximilien/zenika/haproxy/haproxy.stats",
		Type:   "backend",
		Item:   "rtime",
	}
	val, err := haprobe.Value()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("rtime: %fms\n", val)
}
