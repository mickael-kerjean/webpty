package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	assert.Equal(t, ErrNotFound.Error(), "Page Not Found")
	assert.Equal(t, ErrNotAuthorized.Error(), "Not Authorized")
	assert.Equal(t, ErrNotAvailable.Error(), "Not Available")
}
