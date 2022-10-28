package common

import (
	"fmt"
	"net"
)

func GetMachineInfo() string {
	return `{"os":"Ubuntu 21.04","arch":"arm","kernel":"Linux 5.11.0-49-generic","hostname":"pi","ip":"public_ip"}`
}

func GetAddress() []string {
	ips := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				ips = append(ips, fmt.Sprintf(ipnet.IP.String()))
			}
		}
	}
	return ips
}
