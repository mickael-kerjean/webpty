package main

import (
	"net/http"

	. "github.com/mickael-kerjean/webpty/common"
	wctrl "github.com/mickael-kerjean/webpty/ctrl"
	"github.com/mickael-kerjean/webpty/webfleet/ctrl"
	m "github.com/mickael-kerjean/webpty/webfleet/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	Log.Info("Listening on %s", addr)
	http.ListenAndServe(addr, router)
}
