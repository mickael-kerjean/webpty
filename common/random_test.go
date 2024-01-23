package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomGenerator(t *testing.T) {
	b := ""
	for i := 1; i < 250; i++ {
		size := (i % 20) + 1
		b = RandomString(size)
		assert.True(t, len(b) > 0)
		assert.Equal(t, size, len(b))
	}
}
