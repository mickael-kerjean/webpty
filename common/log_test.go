package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	Log.Info("hello %s", "world")
	Log.Warning("hello %s", "world")
	Log.Error("hello %s", "world")
	Log.Debug("hello %s", "world")
	Log.Stdout("hello %s", "world")
	assert.NotNil(t, NewNilLogger())
}
