package handler

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/mickael-kerjean/webpty/common"
	"github.com/rancher/remotedialer"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

var (
	TunnelURL    string
	TunnelServer string = "localhost:8123"
)

func SetupTunnel(res http.ResponseWriter, req *http.Request) {
	tenant := RandomString(5)
	TunnelURL = fmt.Sprintf("http://%s/%s/", TunnelServer, tenant)
	go func() {
		if err := setup(TunnelServer, tenant); err != nil {
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

var Letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		max := *big.NewInt(int64(len(Letters)))
		r, err := rand.Int(rand.Reader, &max)
		if err != nil {
			b[i] = Letters[0]
		} else {
			b[i] = Letters[r.Int64()]
		}
	}
	return string(b)
}
