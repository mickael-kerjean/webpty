package handler

import (
	. "github.com/mickael-kerjean/virtualshell/common"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/ssh"
	"net"
	"net/http"
	"time"
)

func Middleware(fn func(res http.ResponseWriter, req *http.Request)) func(res http.ResponseWriter, req *http.Request) {
	tmpCache := cache.New(5*time.Minute, 10*time.Minute)
	return func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		username, password, ok := req.BasicAuth()
		defer func() {
			Log.Info(
				"HTTP %.1fms %s %s",
				float32(time.Now().Sub(startTime).Microseconds())/1000,
				func() string {
					if ok == false {
						return "anonymous"
					}
					return username
				}(),
				req.URL.Path,
			)
		}()
		if ok == false {
			Log.Error("basic authentication error")
			res.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			ErrorPage(res, ErrNotAuthorized, http.StatusUnauthorized)
			return
		}

		if val, found := tmpCache.Get(username + ":" + password); found == true && val != true {
			client, err := ssh.Dial("tcp", "127.0.0.1:22", &ssh.ClientConfig{
				User: username,
				Auth: []ssh.AuthMethod{ssh.Password(password)},
				HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
					return nil
				},
			})
			if err != nil {
				Log.Error("sshd authentication error: %s", err.Error())
				res.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				ErrorPage(res, ErrNotAuthorized, http.StatusUnauthorized)
				return
			}
			client.Close()
			tmpCache.Set(username+":"+password, true, cache.DefaultExpiration)
		}
		fn(res, req)
	}
}
