package metrics

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNew(t *testing.T) {
	type args struct {
		serviceName    string
		serviceVersion string
	}
	tests := []struct {
		name string
		args args
		want *Prometheus
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.serviceName, tt.args.serviceVersion); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrometheus_MetricsMiddleware(t *testing.T) {
	type fields struct {
		reqCount    *prometheus.CounterVec
		reqLatency  *prometheus.HistogramVec
		reqInFlight *prometheus.GaugeVec
	}
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prometheus{
				reqCount:    tt.fields.reqCount,
				reqLatency:  tt.fields.reqLatency,
				reqInFlight: tt.fields.reqInFlight,
			}
			if got := p.MetricsMiddleware(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prometheus.MetricsMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
