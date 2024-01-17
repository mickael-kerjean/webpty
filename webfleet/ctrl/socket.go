package ctrl

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	. "github.com/mickael-kerjean/webpty/common"
	"github.com/mickael-kerjean/webpty/webfleet/model"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
)

var srv *remotedialer.Server

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})
	// logrus.SetLevel(logrus.DebugLevel)
	// remotedialer.PrintTunnelData = true
	logrus.SetLevel(logrus.FatalLevel)

	authorizer := func(req *http.Request) (string, bool, error) {
		id := req.Header.Get("x-machine-id")
		return id, id != "", nil
	}
	srv = remotedialer.New(
		authorizer,
		func(rw http.ResponseWriter, req *http.Request, code int, err error) { // error writer
			rw.WriteHeader(code)
			rw.Write([]byte(err.Error()))
		},
	)
	srv.PeerToken = "token"
	srv.PeerID = "id"
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ClientSocket(w http.ResponseWriter, r *http.Request, dialer remotedialer.Dialer, url string, tenant string) {
	proxyConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Log.Error("socket.go::upgrade_error tenant[%s] err[%s]", tenant, err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	appConn, _, err := (&websocket.Dialer{
		NetDial:         dialer,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}).Dial(
		url,
		func() http.Header {
			h := http.Header{}
			for k, v := range r.Header {
				switch k {
				case "Sec-Websocket-Key":
				case "Sec-Websocket-Extensions":
				case "Sec-Websocket-Version":
				case "Upgrade":
				case "Connection":
				case "Origin":
				default:
					h.Add(k, strings.Join(v, " "))
				}
			}
			return h
		}(),
	)
	if err != nil {
		Log.Error("socket.go::connection_failed tenant[%s] err[%s]", tenant, err.Error())
		return
	}
	Log.Info("connected to %s", url)
	go func() {
		for {
			mtype, message, err := appConn.ReadMessage()
			if err != nil {
				Log.Error("socket.go::conn::read_message tenant[%s] err[%s]", tenant, err.Error())
				return
			}
			err = proxyConn.WriteMessage(mtype, message)
			if err != nil {
				Log.Error("socket.go::conn::write_message tenant[%s] err[%s]", tenant, err.Error())
				return
			}
		}
	}()
	for {
		mtype, message, err := proxyConn.ReadMessage()
		if err != nil {
			Log.Error("socket.go::proxy::read_message tenant[%s] err[%s]", tenant, err.Error())
			return
		}
		err = appConn.WriteMessage(mtype, message)
		if err != nil {
			Log.Error("socket.go::proxy::write_message tenant[%s] err[%s]", tenant, err.Error())
			return
		}
	}
}

func ClientHTTP(w http.ResponseWriter, r *http.Request, dialer remotedialer.Dialer, url string, tenant string) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error("REQ ERR %s: %s", url, err.Error())
		remotedialer.DefaultErrorWriter(w, r, 500, err)
		return
	}
	for k, v := range r.Header {
		r.Header.Add(k, strings.Join(v, " "))
	}
	resp, err := (&http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Dial:            dialer,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}).Do(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		w.Header().Add(k, strings.Join(v, " "))
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	Log.Info("HTTP tenant=%s url=%s", tenant, url)
}

func TunnelConnect(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("x-machine-id")
	jsonStr := r.Header.Get("x-machine-info")
	data := map[string]interface{}{}
	json.Unmarshal([]byte(jsonStr), &data)
	go func() {
		model.Machines.Add(id, data)
		model.Machines.Vacuum(id, srv.Dialer(id, 15*time.Second))
	}()
	srv.ServeHTTP(w, r)
}

func TunnelMain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if r.Header.Get("Connection") == "Upgrade" {
		ClientSocket(
			w, r,
			srv.Dialer(vars["tenant"], 15*time.Second),
			fmt.Sprintf("wss://localhost:3456/%s", vars["path"]),
			vars["tenant"],
		)
		return
	}
	ClientHTTP(
		w, r,
		srv.Dialer(vars["tenant"], 15*time.Second),
		fmt.Sprintf("https://localhost:3456/%s", vars["path"]),
		vars["tenant"],
	)
}
