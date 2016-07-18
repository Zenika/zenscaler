package probe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandProbe(t *testing.T) {
	assert := assert.New(t)
	cp := NewCommand("echo 0.15")
	assert.Equal(0.15, cp.Value())
}
