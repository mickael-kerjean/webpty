package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	. "github.com/mickael-kerjean/webpty/common"
	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

func authorizer(req *http.Request) (string, bool, error) {
	id := req.Header.Get("x-machine-id")
	return id, id != "", nil
}

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
	}).Dial(url, nil)
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
	resp, err := (&http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Dial:            dialer,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}).Get(url)
	if err != nil {
		Log.Error("REQ ERR %s: %s", url, err.Error())
		remotedialer.DefaultErrorWriter(rw, req, 500, err)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		for _, h := range v {
			if rw.Header().Get(k) == "" {
				rw.Header().Set(k, h)
			} else {
				rw.Header().Add(k, h)
			}
		}
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

	handler := remotedialer.New(authorizer, remotedialer.DefaultErrorWriter)
	handler.PeerToken = "token"
	handler.PeerID = "id"

	router := mux.NewRouter()
	router.Handle("/connect", handler)
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
	fmt.Println("Listening on ", addr)
	http.ListenAndServe(addr, router)
}
