package common

import (
	"encoding/json"
	"github.com/mickael-kerjean/webpty/webfleet/model"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
)

func GetMachineInfo() []byte {
	s := model.Server{
		Hostname: func() string {
			content, err := ioutil.ReadFile("/etc/hostname")
			if err != nil {
				return ""
			}
			return strings.TrimSpace(string(content))
		}(),
		Os: func() string {
			content, err := ioutil.ReadFile("/etc/os-release")
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

func GetAddress() []string {
	ips := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() == nil {
				continue
			}
			ip := ipnet.IP.String()
			ips = append(ips, ip)
		}
	}

	ipsToDisplay := []string{}
	for i := 0; i < len(ips); i++ {
		if strings.HasPrefix(ips[i], "172.") && strings.HasSuffix(ips[i], ".0.1") {
			continue
		}
		ipsToDisplay = append(ipsToDisplay, ips[i])
	}
	if len(ipsToDisplay) == 0 {
		return ips
	}
	return ipsToDisplay
}
