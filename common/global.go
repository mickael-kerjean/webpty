package common

import (
	"os"
)

var (
	FLEET_MODE bool
	FLEET_SRV  string
)

func init() {
	FLEET_SRV = os.Getenv("FLEET")
	if FLEET_SRV != "" {
		FLEET_MODE = true
	}
}
