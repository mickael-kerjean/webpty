package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMachineInfo(t *testing.T) {
	s := GetMachineInfo()

	var data map[string]interface{}
	err := json.Unmarshal(s, &data)
	assert.Nil(t, err)
	assert.NotEqual(t, "", data["hostname"])
	assert.NotEqual(t, "", data["os"])
	assert.NotEqual(t, "", data["kernel"])
	assert.NotEqual(t, "", data["arch"])
	// assert.NotEqual(t, data["publicIP"], "")
	// assert.Equal(t, data["privateIP"], "")
	// Log.Info("DATA %+v", data)
}
