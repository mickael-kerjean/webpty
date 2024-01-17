package model

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

var (
	Machines    = ServerManager{}
	VACUUM_TIME = 5 * time.Second
)

type Server struct {
	Id        string `json:"id"`
	MachineID string `json:"machineID"`
	Device    string `json:"device"`
	Hostname  string `json:"hostname"`
	Os        string `json:"os"`
	Kernel    string `json:"kernel"`
	Arch      string `json:"arch"`
	PublicIP  string `json:"publicIP"`
	PrivateIP string `json:"privateIP"`
	IsOnline  bool   `json:"isOnline"`
}

type ServerManager struct {
	db sync.Map
}

func (this *ServerManager) Add(tenant string, info map[string]interface{}) error {
	this.db.Range(func(key any, value any) bool {
		if value.(Server).MachineID == info["machineID"] {
			this.db.Delete(value.(Server).Id)
		}
		return true
	})
	this.db.Store(tenant, Server{
		tenant, str(info["machineID"]), str(info["device"]),
		str(info["hostname"]), str(info["os"]), str(info["kernel"]), str(info["arch"]),
		str(info["publicIp"]), str(info["privateIP"]), true,
	})
	return nil
}

func (this *ServerManager) Remove(tenant string) {
	this.db.Delete(tenant)
}

func (this *ServerManager) List() ([]Server, error) {
	s := []Server{}
	this.db.Range(func(key any, value any) bool {
		s = append(s, value.(Server))
		return true
	})
	sort.Slice(s, func(i, j int) bool {
		return s[i].Hostname < s[j].Hostname
	})
	return s, nil
}

func (this *ServerManager) Vacuum(tenant string, dialer func(network, address string) (net.Conn, error)) {
	lastSeen := time.Now()
	for {
		time.Sleep(VACUUM_TIME)
		conn, err := dialer("tcp", "127.0.0.1:3456")
		if err != nil {
			if time.Since(lastSeen).Seconds() > 60*30 {
				this.db.Delete(tenant)
				return
			} else if v, ok := this.db.Load(tenant); ok {
				vs := v.(Server)
				if vs.IsOnline {
					vs.IsOnline = false
					this.db.Store(tenant, vs)
				}
			}
			continue
		}
		lastSeen = time.Now()
		conn.Close()
	}
}

func str(in interface{}) string {
	if in == nil {
		return "N/A"
	}
	return fmt.Sprintf("%s", in)
}
