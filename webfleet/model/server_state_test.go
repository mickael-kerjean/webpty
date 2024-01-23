package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerState(t *testing.T) {
	// given
	s := ServerManager{}

	// when
	a, err := s.List()
	assert.NoError(t, err)
	// then
	assert.Equal(t, len(a), 0)

	// when
	err = s.Add("tenant", map[string]interface{}{"device": "test"})
	assert.NoError(t, err)
	a, err = s.List()
	assert.NoError(t, err)
	// then
	assert.Equal(t, len(a), 1)

	// when
	s.Remove("tenant")
	a, err = s.List()
	assert.NoError(t, err)
	// then
	assert.Equal(t, len(a), 0)
}
