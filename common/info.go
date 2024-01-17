package common

import (
	"net"
	"strings"
)

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
