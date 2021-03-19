package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	rand.Seed(time.Now().Unix())

	requestDuration := promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "app_request_duration_seconds",
		Help: "Request duration in seconds",
	})

	requestErrorsCount := promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_request_errors_count",
		Help: "Number of errors observed in requests",
	})

	go func() {
		for {
			// Observe a request that took a random amount of time between [0,
			// 10s), which matches the default buckets defined by Prometheus for
			// histogram metrics.

			requestDuration.Observe(rand.Float64() * 10.0)

			// Simulate the failure of one request in ten.

			if rand.Intn(10) == 0 {
				requestErrorsCount.Inc()
			}

			// Simulate a rate of about 1 req/s.

			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
