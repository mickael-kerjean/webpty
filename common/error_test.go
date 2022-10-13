package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	assert.Contains(t, ErrNotFound.Error(), "Page Not Found")
	assert.Contains(t, ErrNotAuthorized.Error(), "Not Authorized")
}
