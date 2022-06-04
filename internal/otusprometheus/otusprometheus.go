package otusprometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	requestsCountMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dialog_requests_count",
		Help: "Dialog server requests count",
	})
	errorsCountMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dialog_errors_count",
		Help: "Dialog server requests count",
	})
	requestTimeMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "dialog_request_time",
		Help: "Dialog server request time",
	})
)

func IncRequestCount() {
	requestsCountMetric.Inc()
}

func IncErrorsCount() {
	errorsCountMetric.Inc()
}

func AddRequestTimeData(rtime float64) {
	requestTimeMetric.Set(rtime)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
