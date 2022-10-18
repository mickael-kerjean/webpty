package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/mickael-kerjean/webpty/common"
	"github.com/mickael-kerjean/webpty/common/ssl"
	. "github.com/mickael-kerjean/webpty/handler"
	"github.com/rancher/remotedialer"
	"io/ioutil"
	"net"
	"net/http"
	// "time"
)

var port int = 3456

//go:embed .assets/favicon.ico
var IconFavicon []byte

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
	})
	mux.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write(IconFavicon)
	})
	mux.HandleFunc("/", Middleware(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/socket" {
			HandleSocket(res, req)
			return
		} else if req.Method == "GET" {
			HandleStatic(res, req)
			return
		}
		ErrorPage(res, ErrNotFound, 404)
		return
	}))

	msg := `
    ██╗    ██╗███████╗██████╗ ██████╗ ████████╗██╗   ██╗
    ██║    ██║██╔════╝██╔══██╗██╔══██╗╚══██╔══╝╚██╗ ██╔╝
    ██║ █╗ ██║█████╗  ██████╔╝██████╔╝   ██║    ╚████╔╝
    ██║███╗██║██╔══╝  ██╔══██╗██╔═══╝    ██║     ╚██╔╝
    ╚███╔███╔╝███████╗██████╔╝██║        ██║      ██║
     ╚══╝╚══╝ ╚══════╝╚═════╝ ╚═╝        ╚═╝      ╚═╝

    Web Interface:
`
	for _, url := range getAddress() {
		msg += fmt.Sprintf("    - %s\n", url)
	}
	Log.Stdout(msg + "\nLOGS:")
	TLSCert, _, err := ssl.GenerateSelfSigned()
	if err != nil {
		Log.Error("ssl.GenerateSelfSigned %s", err.Error())
		return
	}
	Log.Info("WebPty is ready to go")

	go setupTunnel()
	if err := (&http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			Certificates: []tls.Certificate{TLSCert},
		},
		ErrorLog: NewNilLogger(),
	}).ListenAndServeTLS("", ""); err != nil {
		Log.Error("[https]: listen_serve %v", err)
	}
}

func getAddress() []string {
	ips := []string{}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				ips = append(ips, fmt.Sprintf("https://%s:%d", ipnet.IP.String(), port))
			}
		}
	}
	return ips
}

func setupTunnel() error {
	proxyURL := "ws://localhost:8123/connect"
	tenant := "test"

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
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	session := remotedialer.NewClientSession(
		func(proto, address string) bool { return true },
		ws,
	)
	defer session.Close()
	go func() {
		Log.Info("setting up tunnel to WebPty")
		_, err = session.Serve(ctx)
		result <- err
	}()

	select {
	case <-ctx.Done():
		Log.Info("Proxy done - url[%s] err[%+v]", proxyURL, ctx.Err())
		return ctx.Err()
	case err := <-result:
		rerr, ok := err.(*websocket.CloseError)
		if ok == false {
			Log.Error("Session serve err - url[%s] err[%s]", proxyURL, err.Error())
		} else if rerr.Code == 1006 {
			Log.Info("Proxy has disconnected")
		} else {
			Log.Error("Session serve code[%d] msg[%s]", rerr.Code, rerr.Text)
		}
		return err
	}
}
