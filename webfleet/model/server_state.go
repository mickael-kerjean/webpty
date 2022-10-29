package model

import (
	"fmt"
	"sync"
)

var Machines ServerManager = ServerManager{
	[]Server{}, sync.Mutex{},
}

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
