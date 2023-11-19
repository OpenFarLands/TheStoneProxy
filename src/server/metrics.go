package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var onlineMetric = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "proxy",
	Subsystem: "proxy",
	Name: "players",

})

func addOnline(online int) {
	onlineMetric.Add(float64(online))
}