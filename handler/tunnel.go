package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/mickael-kerjean/webpty/common"
	"github.com/rancher/remotedialer"
	"io/ioutil"
	"net/http"
	"time"
)

var TunnelURL string

func SetupTunnel(res http.ResponseWriter, req *http.Request) {
	tenant := "test"
	proxy := "localhost:8123"

	TunnelURL = fmt.Sprintf("http://%s/%s/", proxy, tenant)
	go func() {
		if err := setup(proxy, tenant); err != nil {
			res.WriteHeader(500)
			res.Write([]byte(err.Error()))
			return
		}
	}()
	time.Sleep(2 * time.Second)
	http.Redirect(res, req, TunnelURL, http.StatusSeeOther)
}

func RedirectTunnel(res http.ResponseWriter, req *http.Request) {
	if TunnelURL == "" {
		res.Write([]byte(""))
		return
	}
	res.Write([]byte(`
    (function() {
        const tunnelURL = "` + TunnelURL + `"; // server generated
        switch(tunnelURL) {
            case "": return;
            case location.href: return;
            default: location.href = tunnelURL;
        }
    })()`))
}

func setup(url string, tenant string) error {
	proxyURL := fmt.Sprintf("ws://%s/connect", url)
	rootCtx := context.Background()
	dialer := &websocket.Dialer{Proxy: http.ProxyFromEnvironment, HandshakeTimeout: remotedialer.HandshakeTimeOut}
	ws, resp, err := dialer.DialContext(
		rootCtx, proxyURL,
		http.Header{"X-Machine-ID": []string{tenant}},
	)
	if err != nil {
		if resp == nil {
			Log.Error("Failed to connect to proxy")
			return err
		} else {
			rb, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				Log.Error("Failed to connect to proxy. Response status: %v - %v. Couldn't read response body (err: %v)", resp.StatusCode, resp.Status, err2)
			} else {
				Log.Error("Failed to connect to proxy. Response status: %v - %v. Response body: %s", resp.StatusCode, resp.Status, rb)
			}
		}
		return err
	}
	defer ws.Close()

	result := make(chan error, 1)
	session := remotedialer.NewClientSession(
		func(proto, address string) bool { return true },
		ws,
	)
	defer session.Close()
	go func() {
		Log.Info("setting up tunnel to WebPty")
		_, err = session.Serve(rootCtx)
		result <- err
	}()
	select {
	case <-rootCtx.Done():
		Log.Info("Proxy done - url[%s] err[%+v]", proxyURL, rootCtx.Err())
	case err := <-result:
		rerr, ok := err.(*websocket.CloseError)
		if ok == false {
			Log.Error("Session serve err - url[%s] err[%s]", proxyURL, err.Error())
		} else if rerr.Code == 1006 {
			Log.Info("Proxy has disconnected")
		} else {
			Log.Error("Session serve code[%d] msg[%s]", rerr.Code, rerr.Text)
		}
	}
	return nil
}
