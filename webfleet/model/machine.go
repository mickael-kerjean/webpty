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
			content := fileString("/etc/machine-id")
			if content != "" {
				return content
			}
			return fileString("/etc/hostname")
		}(),
		Device: func() string {
			if isMac() {
				return "Apple"
			} else if isAndroid() {
				return execCmd("getprop", "ro.product.brand")
			}
			return fileString("/sys/devices/virtual/dmi/id/sys_vendor") + " " +
				fileString("/sys/devices/virtual/dmi/id/product_version") + " " +
				fileString("/sys/devices/virtual/dmi/id/product_name")
		}(),
		Hostname: func() string {
			if isMac() {
				return execCmd("scutil", "--get", "HostName")
			} else if isAndroid() {
				return "android"
			}
			return fileString("/etc/hostname")
		}(),
		Os: func() string {
			if isMac() {
				return execCmd("sw_vers", "-productName") + " " + execCmd("sw_vers", "-productVersion")
			} else if isAndroid() {
				return "Android"
			}
			content := fileString("/etc/os-release")
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
			return execCmd("uname", "-r")
		}(),
		Arch: func() string {
			if isMac() {
				return execCmd("arch")
			} else if isAndroid() {
				return "arm"
			}
			return execCmd("uname", "-m")
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

func execCmd(program string, args ...string) string {
	c, b := exec.Command(program, args...), new(strings.Builder)
	c.Stdout = b
	c.Run()
	return strings.TrimSpace(b.String())
}

func fileString(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(content))
}

func isMac() bool {
	return runtime.GOOS == "darwin"
}

func isAndroid() bool {
	_, err := os.Stat("/system/app")
	return err == nil
}
