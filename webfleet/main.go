package main

import (
	"github.com/gorilla/mux"
	. "github.com/mickael-kerjean/webpty/common"
	wctrl "github.com/mickael-kerjean/webpty/ctrl"
	"github.com/mickael-kerjean/webpty/webfleet/ctrl"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	addr := ":8123"
	router := mux.NewRouter()
	router.HandleFunc("/connect", ctrl.TunnelConnect)
	router.HandleFunc("/{tenant}/{path:.*}", ctrl.TunnelMain)
	router.HandleFunc("/", ctrl.ListServers)
	router.HandleFunc("/favicon.ico", wctrl.ServeFavicon)
	router.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

	Log.Info("Listening on %s", addr)
	http.ListenAndServe(addr, router)
}
