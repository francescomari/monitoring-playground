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
		mu               sync.RWMutex
		maxDuration      = 10
		errorsPercentage = 10
		requestRate      = 1
	)

	log.Printf("using max duration %v", maxDuration)
	log.Printf("using errors percentage %v", errorsPercentage)
	log.Printf("using request rate %v", requestRate)

	go func() {
		for {
			mu.RLock()

			maxDurationValue := maxDuration
			errorsPercentageValue := errorsPercentage
			requestRateValue := requestRate

			mu.RUnlock()

			// Observe a request that took a random amount of time between (0,
			// N) seconds. The default for N is 10s, which fits the highest
			// bucket defined by default by a Prometheus histogram.

			requestDuration.Observe(float64(rand.Intn(maxDurationValue)))

			// Simulate the failure of a certain percentage of the requests.

			if rand.Intn(100) < errorsPercentageValue {
				requestErrorsCount.Inc()
			}

			// Simulate the configured request rate.

			time.Sleep(time.Duration(float64(time.Second) / float64(requestRateValue)))
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

		maxDurationValue, hasMaxDurationValue, err := intValue(r, "maxDuration")
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

		errorsPercentageValue, hasErrorsPercentageValue, err := intValue(r, "errorsPercentage")
		if err != nil {
			log.Printf("error: parse errors percentage: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if hasErrorsPercentageValue && (errorsPercentageValue < 0 || errorsPercentageValue > 100) {
			log.Printf("error: invalid errors percentage '%v'", errorsPercentageValue)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		requestRateValue, hasRequestRateValue, err := intValue(r, "requestRate")
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

		if hasErrorsPercentageValue {
			log.Printf("setting errors percentage to %v", errorsPercentageValue)
			errorsPercentage = errorsPercentageValue
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

func intValue(r *http.Request, name string) (int, bool, error) {
	if _, ok := r.PostForm[name]; !ok {
		return 0, false, nil
	}

	v, err := strconv.ParseInt(r.PostForm.Get(name), 10, 64)
	if err != nil {
		return 0, false, err
	}

	return int(v), true, nil
}
