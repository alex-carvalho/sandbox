package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	healthCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"endpoint", "status"},
	)
)

func init() {
	prometheus.MustRegister(healthCounter)
}

func getAppVersion() string {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		return "1"
	}
	return version
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	healthCounter.WithLabelValues("/health", "200").Inc()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	version := getAppVersion()
	if version == "1" {
		healthCounter.WithLabelValues("/", "200").Inc()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "canary app v1 - OK\n")
	} else {
		healthCounter.WithLabelValues("/", "500").Inc()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "canary app v2 - ERROR\n")
	}
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/health", healthHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting canary exporter on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
