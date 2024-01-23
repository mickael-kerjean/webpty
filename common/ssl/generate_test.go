package ssl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSLGeneration(t *testing.T) {
	cert, pool, err := GenerateSelfSigned()

	assert.NoError(t, err)
	assert.NotNil(t, cert)
	assert.NotNil(t, pool)
}
