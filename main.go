package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
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

	var (
		mu          sync.RWMutex
		maxDuration = 10.0
		errorsRatio = 0.1
		requestRate = 1.0
	)

	log.Printf("using max duration %v", maxDuration)
	log.Printf("using errors ratio %v", errorsRatio)
	log.Printf("using request rate %v", requestRate)

	go func() {
		for {
			mu.RLock()

			maxDurationValue := maxDuration
			errorsRatioValue := errorsRatio
			requestRateValue := requestRate

			mu.RUnlock()

			// Observe a request that took a random amount of time between (0,
			// N) seconds. The default for N is 10s, which fits the highest
			// bucket defined by default by a Prometheus histogram.

			requestDuration.Observe(rand.Float64() * maxDurationValue)

			// Simulate the failure of a certain percentage of the requests.

			if rand.Float64() < errorsRatioValue {
				requestErrorsCount.Inc()
			}

			// Simulate the configured request rate.

			time.Sleep(time.Duration(float64(time.Second) / requestRateValue))
		}
	}()

	http.HandleFunc("/limits", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("error: parse form: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		maxDurationValue, hasMaxDurationValue, err := floatValue(r, "maxDuration")
		if err != nil {
			log.Printf("error: parse max duration: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if hasMaxDurationValue && maxDurationValue <= 0 {
			log.Printf("error: invalid max duration '%v'", maxDurationValue)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		errorsRatioValue, hasErrorsRatioValue, err := floatValue(r, "errorsRatio")
		if err != nil {
			log.Printf("error: parse errors ratio: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if hasErrorsRatioValue && (errorsRatioValue < 0 || errorsRatioValue > 1) {
			log.Printf("error: invalid errors ratio '%v'", errorsRatioValue)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		requestRateValue, hasRequestRateValue, err := floatValue(r, "requestRate")
		if err != nil {
			log.Printf("error: parse request rate value: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if hasRequestRateValue && requestRateValue <= 0 {
			log.Printf("error: invalid request rate value '%v'", requestRateValue)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		mu.Lock()

		if hasMaxDurationValue {
			log.Printf("setting max duration to %v", maxDurationValue)
			maxDuration = maxDurationValue
		}

		if hasErrorsRatioValue {
			log.Printf("setting errors ratio to %v", errorsRatioValue)
			errorsRatio = errorsRatioValue
		}

		if hasRequestRateValue {
			log.Printf("setting request rate to %v", requestRateValue)
			requestRate = requestRateValue
		}

		mu.Unlock()
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func floatValue(r *http.Request, name string) (float64, bool, error) {
	if _, ok := r.PostForm[name]; !ok {
		return 0, false, nil
	}

	v, err := strconv.ParseFloat(r.PostForm.Get(name), 64)
	if err != nil {
		return 0, false, err
	}

	return v, true, nil
}
