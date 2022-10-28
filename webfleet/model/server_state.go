package model

import (
	"fmt"
	"sync"
)

var Machines ServerManager = ServerManager{[]Server{
	// Server{"test", "mickael", "Ubuntu 21.04", "Linux 5.11.0-49-generic", "x86-64", "127.0.0.1", "127.0.0.1", true},
	// Server{"", "awuavadmx01.citc.health.nsw.gov.au", "Oracle Linux Server 7.9", "Linux 5.4.17-2135.311.6.el7uek.x86_64", "x86_64", "127.0.0.1", "127.0.0.1", true},
	// Server{"", "EHL5CG00917SY", "Ubuntu 20.04.1 LTS", "", "x86-64", "127.0.0.1", "127.0.0.1", true},
	// Server{"", "pi", "Raspbian GNU/Linux 10 (buster)", "Kernel: Linux 5.10.103-v7+", "arm", "127.0.0.1", "127.0.0.1", true},
	// Server{"", "rick", "Ubuntu 18.04.5 LTS", "", "x86-64", "127.0.0.1", "127.0.0.1", false},
}, sync.Mutex{}}

type Server struct {
	Id        string
	Hostname  string `json:"hostname"`
	Os        string `json:"os"`
	Kernel    string `json:"kernel"`
	Arch      string `json:"arch"`
	PublicIP  string `json:"publicIp"`
	PrivateIP string `json:"privateIp"`
	IsOnline  bool   `json:"isOnline"`
}

type ServerManager struct {
	db []Server
	sync.Mutex
}

func (this *ServerManager) Add(key string, info map[string]interface{}) error {
	this.Lock()
	this.db = append(this.db, Server{
		key,
		str(info["hostname"]), str(info["os"]), str(info["kernel"]), str(info["arch"]),
		str(info["publicIp"]), str(info["privateIP"]), true,
	})
	this.Unlock()
	return nil
}

func (this ServerManager) Remove(key string) {
	// TODO
}

func (this ServerManager) List() ([]Server, error) {
	return this.db, nil
}

func str(in interface{}) string {
	if in == nil {
		return "N/A"
	}
	return fmt.Sprintf("%s", in)
}
