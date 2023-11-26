package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var onlineMetric = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "stone",
	Subsystem: "proxy",
	Name: "players",
	Help: "Total number of online players right now.",
})

func addOnline(online int) {
	onlineMetric.Add(float64(online))
}