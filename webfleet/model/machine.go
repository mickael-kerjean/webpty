package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	. "github.com/mickael-kerjean/webpty/common"
)

func GetMachineInfo() []byte {
	s := Server{
		MachineID: func() string {
			if runtime.GOOS == "darwin" {
				c, b := exec.Command("system_profiler", "SPHardwareDataType"), new(strings.Builder)
				c.Stdout = b
				c.Run()
				return createHash(b.String())
			}
			content, err := os.ReadFile("/etc/machine-id")
			if err != nil {
				return ""
			}
			return strings.TrimSpace(string(content))
		}(),
		Device: func() string {
			if runtime.GOOS == "darwin" {
				return "Apple"
			}
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
			if runtime.GOOS == "darwin" {
				c, b := exec.Command("scutil", "--get", "HostName"), new(strings.Builder)
				c.Stdout = b
				c.Run()
				return b.String()
			}
			content, err := os.ReadFile("/etc/hostname")
			if err != nil {
				return ""
			}
			return strings.TrimSpace(string(content))
		}(),
		Os: func() string {
			if runtime.GOOS == "darwin" {
				a, b := exec.Command("sw_vers", "-productName"), new(strings.Builder)
				a.Stdout = b
				a.Run()
				c, d := exec.Command("sw_vers", "-productVersion"), new(strings.Builder)
				c.Stdout = d
				c.Run()
				return b.String() + " " + d.String()
			}
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
			if runtime.GOOS == "darwin" {
				return ""
			}
			c, b := exec.Command("uname", "-r"), new(strings.Builder)
			c.Stdout = b
			c.Run()
			return strings.TrimSpace(b.String())
		}(),
		Arch: func() string {
			if runtime.GOOS == "darwin" {
				c, b := exec.Command("arch"), new(strings.Builder)
				c.Stdout = b
				c.Run()
				return strings.TrimSpace(b.String())
			}
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

func createHash(text string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, text)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
