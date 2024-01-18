package main

import (
	"crypto/tls"
	"net/http"
	"os"

	. "github.com/mickael-kerjean/webpty/common"
	wctrl "github.com/mickael-kerjean/webpty/ctrl"
	"github.com/mickael-kerjean/webpty/webfleet/ctrl"
	m "github.com/mickael-kerjean/webpty/webfleet/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	addr := ":8123"
	router := mux.NewRouter()
	router.HandleFunc("/connect", ctrl.TunnelConnect)
	router.HandleFunc("/{tenant}/healthz", ctrl.TunnelMain)
	router.HandleFunc("/{tenant}/{path:.*}", m.WithAuth(ctrl.TunnelMain))
	router.HandleFunc("/", m.WithAuth(ctrl.ListServers))
	router.HandleFunc("/favicon.ico", wctrl.ServeFavicon)
	router.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

	if host := os.Getenv("CERTBOT"); host != "" {
		certManager := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(host, host+addr),
			Cache:      autocert.DirCache("certs"),
		}
		srv := &http.Server{
			Addr:    ":https",
			Handler: router,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
				MinVersion:     tls.VersionTLS12,
			},
		}
		Log.Info("Listening on %s", host)
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		if err := srv.ListenAndServeTLS("", ""); err != nil {
			Log.Error("finalise with %s", err.Error())
		}
		return
	}
	Log.Info("Listening on %s", addr)
	http.ListenAndServe(addr, router)
}
