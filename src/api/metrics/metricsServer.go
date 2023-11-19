package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	conf "github.com/OpenFarLands/TheStoneProxy/src/config"
)

func Setup(paramConfig *conf.Config) error {
	log.Printf("Starting prometheus server on %v.", paramConfig.PrometheusAddress)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(paramConfig.PrometheusAddress, mux)
}
