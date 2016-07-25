package probe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandProbe(t *testing.T) {
	assert := assert.New(t)
	cp := &Command{
		Cmd: "echo 0.15",
	}
	val, err := cp.Value()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(0.15, val)
}
