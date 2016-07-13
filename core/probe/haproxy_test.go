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
	fmt.Printf("rtime: %fms\n", haprobe.Value())
}
