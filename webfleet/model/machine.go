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
			if isMac() {
				return createHash(execCmd("system_profiler", "SPHardwareDataType"))
			} else if isAndroid() {
				return createHash(execCmd("getprop", "ro.product.model") +
					execCmd("getprop", "ro.product.brand") +
					execCmd("getprop", "ro.product.name") +
					execCmd("getprop", "ro.product.manufacturer") +
					execCmd("getprop", "ro.build.version.release") +
					execCmd("getprop", "ro.build.version.sdk") +
					execCmd("getprop", "ro.hardware") +
					execCmd("getprop", "ro.arch") +
					execCmd("getprop", "ro.sf.lcd_density") +
					execCmd("getprop", "persist.sys.display.size") +
					execCmd("getprop", "gsm.operator.alpha") +
					execCmd("getprop", "gsm.version.ril-impl") +
					execCmd("getprop", "getprop bluetooth.nam") +
					execCmd("getprop", "wifi.interface") +
					execCmd("getprop", "ro.serialno") +
					execCmd("getprop", "ro.bootloader") +
					execCmd("getprop", "ro.serialno") +
					execCmd("getprop", "ro.kernel.android.checkjni") +
					execCmd("getprop", "ro.boot.kernel") +
					execCmd("getprop", "persist.sys.country") +
					execCmd("getprop", "persist.sys.language"))
			}
			content, err := os.ReadFile("/etc/machine-id")
			if err != nil {
				return ""
			} else if string(content) != "" {
				return strings.TrimSpace(string(content))
			}
			return RandomString(5)
		}(),
		Device: func() string {
			if isMac() {
				return "Apple"
			} else if isAndroid() {
				return execCmd("getprop", "ro.product.brand")
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
			if isMac() {
				return execCmd("scutil", "--get", "HostName")
			} else if isAndroid() {
				return "android"
			}
			content, err := os.ReadFile("/etc/hostname")
			if err != nil {
				return ""
			}
			return strings.TrimSpace(string(content))
		}(),
		Os: func() string {
			if isMac() {
				return execCmd("sw_vers", "-productName") + " " + execCmd("sw_vers", "-productVersion")
			} else if isAndroid() {
				return "Android"
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
			if isMac() {
				return ""
			} else if isAndroid() {
				return ""
			}
			c, b := exec.Command("uname", "-r"), new(strings.Builder)
			c.Stdout = b
			c.Run()
			return strings.TrimSpace(b.String())
		}(),
		Arch: func() string {
			if isMac() {
				c, b := exec.Command("arch"), new(strings.Builder)
				c.Stdout = b
				c.Run()
				return strings.TrimSpace(b.String())
			} else if isAndroid() {
				return "arm"
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

func isMac() bool {
	return runtime.GOOS == "darwin"
}

func isAndroid() bool {
	_, err := os.Stat("/system/app")
	return err == nil
}

func execCmd(program string, args ...string) string {
	c, b := exec.Command(program, args...), new(strings.Builder)
	c.Stdout = b
	c.Run()
	return b.String()
}
