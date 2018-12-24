package main

import (
	"log"
	"net/http"
	"time"

	// TODO: Move prometheus related things to exporter/prometheus.go
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/waltsu/pgpool2-prometheus-exporter/exporter"
)

var (
	nodeCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_count",
		Help: "How many nodes are in the pool at the moment",
	})
)

func main() {
	go startMetricGathering()
	startPrometheusServer()
}

func startPrometheusServer() {
	log.Println("Starting Prometheus HTTP server")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startMetricGathering() {
	log.Println("Start gathering metrics")

	registerPrometheusMetrics()

	commandExecutor := new(exporter.BashExecutor)
	//commandExecutor := new(exporter.TestExecutor)
	pgpool := exporter.NewPgPool(commandExecutor)

	for {
		gatherMetrics(pgpool)
		time.Sleep(1 * time.Second)
	}
}

func registerPrometheusMetrics() {
	prometheus.MustRegister(nodeCountGauge)
}

func gatherMetrics(pgpool *exporter.PgPool) {
	nodeCount, err := pgpool.GetNodeCount()

	if err != nil {
		log.Println(err)
		return
	}
	nodeCountGauge.Set(float64(nodeCount))
}
