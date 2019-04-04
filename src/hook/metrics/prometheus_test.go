package metrics

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func Test_newLoggingResponseWriter(t *testing.T) {
	type args struct {
		w http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
		want *statusCodeResponseWriter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLoggingResponseWriter(tt.args.w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLoggingResponseWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_statusCodeResponseWriter_WriteHeader(t *testing.T) {
	type fields struct {
		ResponseWriter http.ResponseWriter
		statusCode     int
	}
	type args struct {
		code int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &statusCodeResponseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				statusCode:     tt.fields.statusCode,
			}
			s.WriteHeader(tt.args.code)
		})
	}
}

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

func TestPrometheus_HandleFunc(t *testing.T) {
	type fields struct {
		reqCount    *prometheus.CounterVec
		reqLatency  *prometheus.HistogramVec
		reqInFlight *prometheus.GaugeVec
	}
	type args struct {
		path string
		next http.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  http.HandlerFunc
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
			got, got1 := p.HandleFunc(tt.args.path, tt.args.next)
			if got != tt.want {
				t.Errorf("Prometheus.HandleFunc() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Prometheus.HandleFunc() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
