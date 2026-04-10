package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"result"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter)
	requestCounter.WithLabelValues("success").Add(0)
	requestCounter.WithLabelValues("error").Add(0)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		requestCounter.WithLabelValues("error").Inc()
	} else {
		requestCounter.WithLabelValues("success").Inc()
	}
	// return ok to not fail on health check, but count errors in Prometheus
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "DB_URL = %s - %s\n", dbUrl, os.Getenv("POD_NAME"))
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting canary exporter on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
