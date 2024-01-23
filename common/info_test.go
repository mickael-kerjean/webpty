package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAddress(t *testing.T) {
	assert.True(t, len(GetAddress()) > 0)
}
