package model

var Machines ServerManager = ServerManager{[]Server{
	Server{"test", "mickael", "Ubuntu 21.04", "Linux 5.11.0-49-generic", "x86-64", "127.0.0.1", "127.0.0.1", true},
	Server{"", "awuavadmx01.citc.health.nsw.gov.au", "Oracle Linux Server 7.9", "Linux 5.4.17-2135.311.6.el7uek.x86_64", "x86_64", "127.0.0.1", "127.0.0.1", true},
	Server{"", "EHL5CG00917SY", "Ubuntu 20.04.1 LTS", "", "x86-64", "127.0.0.1", "127.0.0.1", true},
	Server{"", "pi", "Raspbian GNU/Linux 10 (buster)", "Kernel: Linux 5.10.103-v7+", "arm", "127.0.0.1", "127.0.0.1", true},
	Server{"", "rick", "Ubuntu 18.04.5 LTS", "", "x86-64", "127.0.0.1", "127.0.0.1", false},
}}

type Server struct {
	Id        string
	Hostname  string `json:"hostname"`
	Os        string `json:"os"`
	Kernel    string `json:"kernel"`
	Arch      string `json:"arch"`
	PublicIp  string `json:"publicIp"`
	PrivateIP string `json:"privateIp"`
	IsOnline  bool   `json:"isOnline"`
}

type ServerManager struct {
	db []Server
}

func (this ServerManager) Add(key string, ip string) {
	// TODO
}

func (this ServerManager) Remove(key string) {
	// TODO
}

func (this ServerManager) List() []Server {
	return this.db
}
