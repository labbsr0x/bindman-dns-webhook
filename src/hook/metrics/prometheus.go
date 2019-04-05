package metrics

import (
	"github.com/go-errors/errors"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// workaround to get status code on middleware
type statusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *statusCodeResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &statusCodeResponseWriter{w, http.StatusOK}
}

func (s *statusCodeResponseWriter) WriteHeader(code int) {
	s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}

type Prometheus struct {
	reqCount    *prometheus.CounterVec
	reqLatency  *prometheus.HistogramVec
	reqInFlight *prometheus.GaugeVec
}

func New(serviceVersion string) *Prometheus {
	if strings.TrimSpace(serviceVersion) == "" {
		panic(errors.New("Service version must be a non-empty string"))
	}

	buildInfoGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "build_info",
			Help:        "Information about build",
			ConstLabels: prometheus.Labels{"version": serviceVersion},
		})
	buildInfoGauge.Set(1)

	prometheus.MustRegister(buildInfoGauge)

	p := &Prometheus{}
	p.reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(p.reqCount)

	p.reqLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "How long it took to process the request, partitioned by status code, method and HTTP path.",
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(p.reqLatency)

	p.reqInFlight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "How many requests are being processed, partitioned method and HTTP path.",
	},
		[]string{"method", "path"},
	)
	prometheus.MustRegister(p.reqInFlight)

	return p
}

func (p *Prometheus) HandleFunc(path string, next http.HandlerFunc) (string, http.HandlerFunc) {
	return path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseWriter := newLoggingResponseWriter(w)

		p.reqInFlight.WithLabelValues(r.Method, path).Inc()

		start := time.Now()
		next.ServeHTTP(responseWriter, r)
		duration := time.Since(start)

		p.reqInFlight.WithLabelValues(r.Method, path).Dec()

		strStatusCode := strconv.Itoa(responseWriter.statusCode)
		p.reqCount.WithLabelValues(strStatusCode, r.Method, path).Inc()
		p.reqLatency.WithLabelValues(strStatusCode, r.Method, path).Observe(duration.Seconds())

	})
}
