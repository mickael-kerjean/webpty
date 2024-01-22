package ctrl

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/mickael-kerjean/webpty/common"
	"github.com/mickael-kerjean/webpty/webfleet/model"

	"github.com/gorilla/websocket"
	"github.com/rancher/remotedialer"
)

var (
	tunnelURL  string
	tunnelDate time.Time
)

func InitTunnel(tunnelServer string) (string, error) {
	tenant := RandomString(5)
	tunnelURL = fmt.Sprintf("https://%s/%s/", tunnelServer, tenant)
	go func() {
		if err := setup(tunnelServer, tenant, model.GetMachineInfo(), 0); err != nil {
			Log.Error("setup tunnel failed %s", err.Error())
			return
		}
	}()
	for i := 0; i < 30; i++ {
		time.Sleep(time.Duration(i*1000+10) * time.Millisecond)
		resp, err := http.Get(tunnelURL + "healthz")
		if err != nil {
			continue
		} else if resp.StatusCode != 200 {
			continue
		}
		Log.Debug("tunnel established")
		return tunnelServer, nil
	}
	return "", ErrNotFound
}

func RedirectTunnel(res http.ResponseWriter, req *http.Request) {
	if tunnelURL == "" {
		res.Write([]byte(""))
		return
	}
	res.Write([]byte(`
    (function() {
        const tunnelURL = "` + tunnelURL + `"; // server generated
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
	proxyURL := fmt.Sprintf("wss://%s/connect", url)
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
			Log.Error("Failed to connect to proxy. Reconnecting ....")
			time.Sleep(time.Duration(retry*5) * time.Second)
			setup(url, tenant, jsonInfo, retry+1)
			return err
		} else {
			rb, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				Log.Error("Failed to connect to proxy '%s'. Response status: %v - %v. Couldn't read response body (err: %v)", proxyURL, resp.StatusCode, resp.Status, err2)
			} else {
				Log.Error("Failed to connect to proxy '%s'. Response status: %v. Response body: %s", proxyURL, resp.StatusCode, rb)
			}
		}
		return err
	}

	result := make(chan error, 1)
	session := remotedialer.NewClientSession(
		func(proto, address string) bool { return true },
		ws,
	)
	go func() {
		if retry == 0 {
			Log.Debug("setting up tunnel to WebPty")
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
			session.Close()
			ws.Close()
			setup(url, tenant, jsonInfo, retry+1)
		} else {
			Log.Error("Session serve code[%d] msg[%s]", rerr.Code, rerr.Text)
		}
	}
	session.Close()
	ws.Close()
	return nil
}
