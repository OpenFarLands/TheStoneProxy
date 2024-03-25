//go:build metrics

package metrics

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	conf "github.com/OpenFarLands/TheStoneProxy/src/config"
)

var config *conf.Config

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		authToken := config.Metrics.PrometheusBearerAuthToken

		// A ⋀ (B ⋁ С)  <=>  (A ⋀ B) ⋁ (A ⋀ C) xDDD
		if authToken != "" && (header == "" || header != fmt.Sprintf("Bearer %v", authToken)) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Setup(paramConfig *conf.Config) error {
	config = paramConfig

	log.Printf("Starting prometheus server on %v.", config.Metrics.PrometheusAddress)

	mux := http.NewServeMux()
	mux.Handle("/metrics", authMiddleware(promhttp.Handler()))

	return http.ListenAndServe(paramConfig.Metrics.PrometheusAddress, mux)
}
