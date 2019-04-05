package metrics

import (
	"net/http"
	"net/http/httptest"
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
		{
			"default status code",
			args{},
			&statusCodeResponseWriter{nil, http.StatusOK},
		},
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
		{
			"statusCodeResponseWriter.statusCode must be the same passed to WriteHeader",
			fields{ResponseWriter: httptest.NewRecorder()},
			args{http.StatusInternalServerError},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &statusCodeResponseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				statusCode:     tt.fields.statusCode,
			}
			s.WriteHeader(tt.args.code)
			if s.statusCode != tt.args.code {
				t.Errorf("want %d, got %d", tt.args.code, s.statusCode)
			}
		})
	}
}

func TestNew(t *testing.T) {
	newMetrics := New("1")
	// validate all metrics have been instantiated
	if newMetrics.reqCount == nil {
		t.Fatalf("expected reqCount to be already instantiated")
	}
	if newMetrics.reqLatency == nil {
		t.Fatalf("expected reqLatency to be already instantiated")
	}
	if newMetrics.reqInFlight == nil {
		t.Fatalf("expected reqInFlight to be already instantiated")
	}

	registerValidation := func(c prometheus.Collector, t *testing.T) {
		err := prometheus.Register(c)
		if err == nil {
			t.Fatal("expected s to be already registered")
		}
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			t.Fatal("unexpected registration error:", err)
		}
	}
	// must have all metrics already registered
	registerValidation(newMetrics.reqCount, t)
	registerValidation(newMetrics.reqLatency, t)
	registerValidation(newMetrics.reqInFlight, t)
}

func TestNew_ValidateArgs(t *testing.T) {
	type args struct {
		serviceVersion string
	}
	type expected struct {
		error bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			"valid service version",
			args{"1"},
			expected{error: false},
		}, {
			"empty service version",
			args{""},
			expected{error: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRegistry()
			defer func() {
				err := recover()
				if tt.expected.error != (err != nil) {
					t.Fatalf("expected an error %t but got %v", tt.expected.error, err)
				}
			}()
			New(tt.args.serviceVersion)
		})
	}
}

// TODO check metrics values after handle execution
func TestPrometheus_HandleFunc(t *testing.T) {
	type want struct {
		path       string
		method     string
		statusCode int
		body       string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			"returned path must be the same received via function parameter",
			want{"/api/get", "GET", http.StatusPaymentRequired, "hello"},
		},
		{
			"returned path must be the same received via function parameter",
			want{"/api/post", "POST", http.StatusOK, ""},
		},
	}

	resetRegistry()
	p := New("1")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// handle to collect metrics from
			next := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.want.statusCode)
				_, err := w.Write([]byte(tt.want.body))
				if err != nil {
					t.Fatal(err)
				}
			}
			path, handle := p.HandleFunc(tt.want.path, next)
			if path != tt.want.path {
				t.Errorf("path must be the same provided by function parameter. got = %v, want %v", path, tt.want.path)
			}

			// metrics middleware uses only method attribute from http request
			req := &http.Request{
				Method: "GET",
			}
			resp := httptest.NewRecorder()

			handle.ServeHTTP(resp, req)

			// metrics middleware cannot change response status code
			if resp.Code != tt.want.statusCode {
				t.Fatalf("expected status %d, got %d", tt.want.statusCode, resp.Code)
			}
			// metrics middleware cannot change response body
			if resp.Body.String() != tt.want.body {
				t.Fatalf("expected body %s, got %s", tt.want.body, resp.Body.String())
			}
		})
	}
}

func resetRegistry() {
	// reset registry
	// Create and assign a new Prometheus Registerer/Gatherer for each test
	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry
	prometheus.DefaultGatherer = registry
}
