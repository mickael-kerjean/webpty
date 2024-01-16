package ctrl

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	. "github.com/mickael-kerjean/webpty/common"

	"github.com/gorilla/websocket"
	"github.com/rancher/remotedialer"
)

var (
	TunnelURL  string
	TunnelDate time.Time
)

func SetupTunnel(res http.ResponseWriter, req *http.Request) {
	tenant := RandomString(5)
	tunnelServer := req.URL.Query().Get("srv")
	TunnelURL = fmt.Sprintf("http://%s/%s/", tunnelServer, tenant)
	go func() {
		if err := setup(tunnelServer, tenant, GetMachineInfo(), 0); err != nil {
			res.WriteHeader(500)
			res.Write([]byte(err.Error()))
			return
		}
	}()
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(TunnelURL + "healthz")
		if err != nil {
			continue
		} else if resp.StatusCode != 200 {
			continue
		}
		http.Redirect(res, req, TunnelURL, http.StatusSeeOther)
		return
	}
	ErrorPage(res, ErrNotFound, 404)
	return
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

func setup(url string, tenant string, jsonInfo []byte, retry int) error {
	if retry > 100 {
		return ErrNotAvailable
	}
	proxyURL := fmt.Sprintf("ws://%s/connect", url)
	rootCtx := context.Background()
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: remotedialer.HandshakeTimeOut,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}
	ws, resp, err := dialer.DialContext(
		rootCtx, proxyURL,
		http.Header{
			"X-Machine-ID":   []string{tenant},
			"X-Machine-Info": []string{string(jsonInfo)},
		},
	)
	if err != nil {
		if resp == nil {
			Log.Error("Failed to connect to proxy. Reconnecting proxy[%s] err[%s] ....", proxyURL, err.Error())
			time.Sleep(time.Duration(retry*5) * time.Second)
			setup(url, tenant, jsonInfo, retry+1)
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
		if retry == 0 {
			Log.Info("setting up tunnel to WebPty")
		} else {
			Log.Info("tunnel is back online")
		}
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
			Log.Info("Proxy has disconnected. Reconnecting ....")
			time.Sleep(time.Duration(retry*5) * time.Second)
			setup(url, tenant, jsonInfo, retry+1)
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
