package model

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	. "github.com/mickael-kerjean/webpty/common"
)

func GetMachineInfo() []byte {
	s := Server{
		MachineID: func() string {
			content, err := os.ReadFile("/etc/machine-id")
			if err != nil {
				return ""
			}
			return strings.TrimSpace(string(content))
		}(),
		Device: func() string {
			m := ""
			if content, err := os.ReadFile("/sys/devices/virtual/dmi/id/sys_vendor"); err == nil && string(content) != "" {
				m += string(content)
			}
			if content, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_version"); err == nil && string(content) != "" {
				m += string(content)
			}
			if content, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_name"); err == nil && string(content) != "" {
				m += string(content)
			}
			return m
		}(),
		Hostname: func() string {
			content, err := os.ReadFile("/etc/hostname")
			if err != nil {
				return ""
			}
			return strings.TrimSpace(string(content))
		}(),
		Os: func() string {
			content, err := os.ReadFile("/etc/os-release")
			if err != nil {
				return ""
			}
			for _, line := range strings.Split(string(content), "\n") {
				if strings.HasPrefix(line, "PRETTY_NAME=") == false {
					continue
				}
				return strings.TrimSuffix(strings.TrimPrefix(line, "PRETTY_NAME=\""), "\"")
			}
			return ""
		}(),
		Kernel: func() string {
			c, b := exec.Command("uname", "-r"), new(strings.Builder)
			c.Stdout = b
			c.Run()
			return strings.TrimSpace(b.String())
		}(),
		Arch: func() string {
			c, b := exec.Command("uname", "-m"), new(strings.Builder)
			c.Stdout = b
			c.Run()
			return strings.TrimSpace(b.String())
		}(),
		PublicIP: func() string {
			return ""
		}(),
		PrivateIP: func() string {
			return ""
		}(),
		IsOnline: func() bool {
			return true
		}(),
	}
	b, err := json.Marshal(s)
	if err != nil {
		Log.Error("common::info machine marshall %s", err.Error())
		return []byte("{}")
	}
	return b
}
