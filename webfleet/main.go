package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	. "github.com/mickael-kerjean/webpty/common"
	wctrl "github.com/mickael-kerjean/webpty/ctrl"
	"github.com/mickael-kerjean/webpty/webfleet/ctrl"
	"github.com/mickael-kerjean/webpty/webfleet/model"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

//go:embed tmpl/*
var tmplFs embed.FS

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ClientSocket(rw http.ResponseWriter, req *http.Request, dialer remotedialer.Dialer, url string) {
	proxyConn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		Log.Error("WS UPGRADE ERR %s", err.Error())
		rw.Write([]byte(err.Error()))
		return
	}
	appConn, _, err := (&websocket.Dialer{
		NetDial:         dialer,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}).Dial(
		url,
		func() http.Header {
			h := http.Header{}
			for k, v := range req.Header {
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
		Log.Error("WS DIAL ERR %s", err.Error())
		return
	}
	go func() {
		for {
			mtype, message, err := appConn.ReadMessage()
			if err != nil {
				Log.Error("WS(read) ReadMessage ERR: %v", err.Error())
				return
			}
			err = proxyConn.WriteMessage(mtype, message)
			if err != nil {
				Log.Error("WS(read) WriteMessage ERR: %v", err.Error())
				return
			}
			Log.Info("WS OUT %s", url)
		}
	}()
	for {
		mtype, message, err := proxyConn.ReadMessage()
		if err != nil {
			Log.Error("WS(write) ReadMessage ERR: %s", err.Error())
			return
		}
		err = appConn.WriteMessage(mtype, message)
		if err != nil {
			Log.Error("WS(write) WriteMessage ERR: %s", err.Error())
			return
		}
		Log.Info("WS IN %s", url)
	}
	return
}

func ClientHTTP(rw http.ResponseWriter, req *http.Request, dialer remotedialer.Dialer, url string) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error("REQ ERR %s: %s", url, err.Error())
		remotedialer.DefaultErrorWriter(rw, req, 500, err)
		return
	}
	for k, v := range req.Header {
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
		Log.Error("REQ ERR %s: %s", url, err.Error())
		remotedialer.DefaultErrorWriter(rw, req, 500, err)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		rw.Header().Add(k, strings.Join(v, " "))
	}
	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
	Log.Info("HTTP OK %s", url)
}

func main() {
	addr := ":8123"
	// remotedialer.PrintTunnelData = true
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		// FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	authorizer := func(req *http.Request) (string, bool, error) {
		id := req.Header.Get("x-machine-id")
		return id, id != "", nil
	}
	handler := remotedialer.New(
		authorizer,
		func(rw http.ResponseWriter, req *http.Request, code int, err error) { // error writer
			rw.WriteHeader(code)
			rw.Write([]byte(err.Error()))
		},
	)
	handler.PeerToken = "token"
	handler.PeerID = "id"

	router := mux.NewRouter()
	router.HandleFunc("/connect", func(rw http.ResponseWriter, req *http.Request) {
		if id, ok, err := authorizer(req); err == nil && ok == true {
			model.Machines.Add(id, req.RemoteAddr)
		}
		handler.ServeHTTP(rw, req)
	})
	router.HandleFunc("/", ctrl.ServeFile(tmplFs, "tmpl/dashboard.html"))
	router.HandleFunc("/{tenant}/{path:.*}", func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		if req.Header.Get("Connection") == "Upgrade" {
			ClientSocket(
				rw, req,
				handler.Dialer(vars["tenant"], 15*time.Second),
				fmt.Sprintf("wss://localhost:3456/%s", vars["path"]),
			)
			return
		}
		ClientHTTP(
			rw, req,
			handler.Dialer(vars["tenant"], 15*time.Second),
			fmt.Sprintf("https://localhost:3456/%s", vars["path"]),
		)
	})
	router.HandleFunc("/favicon.ico", wctrl.ServeFavicon)
	router.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

	fmt.Println("Listening on ", addr)
	http.ListenAndServe(addr, router)
}
