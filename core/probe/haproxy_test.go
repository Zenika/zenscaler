package probe

import (
	"fmt"
	"testing"
)

func TestHAproxyRTimeProbe(t *testing.T) {
	haprobe := &HAproxy{
		Socket: "/home/maximilien/zenika/haproxy/haproxy.stats",
		Item:   "req_rate",
	}
	data, _ := haprobe.getStats("backend")
	fmt.Printf("rtime: %sms\n", data["rtime"][1])
}
