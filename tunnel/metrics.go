package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rancher/remotedialer/metrics"
)

func init() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(metrics.TotalAddWS)
	reg.MustRegister(metrics.TotalRemoveWS)
	reg.MustRegister(metrics.TotalAddConnectionsForWS)
	reg.MustRegister(metrics.TotalRemoveConnectionsForWS)
	reg.MustRegister(metrics.TotalTransmitBytesOnWS)
	reg.MustRegister(metrics.TotalTransmitErrorBytesOnWS)
	reg.MustRegister(metrics.TotalReceiveBytesOnWS)
	reg.MustRegister(metrics.TotalAddPeerAttempt)
	reg.MustRegister(metrics.TotalPeerConnected)
	reg.MustRegister(metrics.TotalPeerDisConnected)

	prometheus.DefaultRegisterer = reg
	prometheus.DefaultGatherer = reg
}
