package handler

import (
	. "github.com/mickael-kerjean/virtualshell/common"
	"golang.org/x/crypto/ssh"
	"net"
	"net/http"
	"time"
)

func Middleware(fn func(res http.ResponseWriter, req *http.Request)) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		username, password, ok := req.BasicAuth()
		defer func() {
			Log.Info(
				"HTTP %dms %s %s",
				time.Now().Sub(startTime).Milliseconds(),
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
		fn(res, req)
	}
}
