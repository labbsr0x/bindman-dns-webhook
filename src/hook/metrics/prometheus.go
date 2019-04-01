package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

type Prometheus struct {
	reqCount    *prometheus.CounterVec
	reqLatency  *prometheus.HistogramVec
	reqInFlight *prometheus.GaugeVec
}

func New(serviceName, serviceVersion string) *Prometheus {
	p := &Prometheus{}
	constLabels := prometheus.Labels{"service": serviceName, "service_version": serviceVersion}
	p.reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "http_requests_total",
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: constLabels,
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(p.reqCount)

	p.reqLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "http_request_duration_seconds",
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: constLabels,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(p.reqLatency)

	p.reqInFlight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "http_requests_in_flight_total",
		Help:        "How many requests are being processed, partitioned method and HTTP path.",
		ConstLabels: constLabels,
	},
		[]string{"method", "path"},
	)
	prometheus.MustRegister(p.reqInFlight)
	return p
}

func (p *Prometheus) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path

		p.reqInFlight.WithLabelValues(r.Method, path).Inc()

		next.ServeHTTP(w, r)

		p.reqInFlight.WithLabelValues(r.Method, path).Dec()

		statusCode := w.Header().Get("Status-Code")

		p.reqCount.WithLabelValues(statusCode, r.Method, path).
			Inc()
		p.reqLatency.WithLabelValues(statusCode, r.Method, path).
			Observe(float64(time.Since(start).Seconds()) / 1000000000)

	})
}
