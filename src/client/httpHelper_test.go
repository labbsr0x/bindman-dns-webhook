package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type args struct {
	method      string
	contentType string
	payload     []byte
}
type expected struct {
	code int
	body []byte
	err  error
}

func TestBindmanHTTPHelper_Put(t *testing.T) {
	ts := createTestServer(t, args{"PUT", "application/json", nil}, expected{http.StatusOK, nil, nil})
	defer ts.Close()
	bhh := &BindmanHTTPHelper{}
	_, _, _ = bhh.Put(ts.URL, nil)
}

func TestBindmanHTTPHelper_Post(t *testing.T) {
	ts := createTestServer(t, args{"POST", "application/json", nil}, expected{http.StatusOK, nil, nil})
	defer ts.Close()
	bhh := &BindmanHTTPHelper{}
	_, _, _ = bhh.Post(ts.URL, nil)
}

func TestBindmanHTTPHelper_Get(t *testing.T) {
	ts := createTestServer(t, args{"GET", "", nil}, expected{http.StatusOK, nil, nil})
	defer ts.Close()
	bhh := &BindmanHTTPHelper{}
	_, _, _ = bhh.Get(ts.URL)
}

func TestBindmanHTTPHelper_Delete(t *testing.T) {
	ts := createTestServer(t, args{"DELETE", "", nil}, expected{http.StatusOK, nil, nil})
	defer ts.Close()
	bhh := &BindmanHTTPHelper{}
	_, _, _ = bhh.Delete(ts.URL)
}

func TestBindmanHTTPHelper_request(t *testing.T) {
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			"happy path",
			args{"PUT", "application/json", []byte(`{"hello": "world"}`)},
			expected{http.StatusOK, []byte(`{"world": "hello"}`), nil},
		},
		{
			"nil request data",
			args{"POST", "application/json", nil},
			expected{http.StatusOK, []byte(`{"world": "hello"}`), nil},
		},
		{
			"nil response body",
			args{"DELETE", "text/plain", nil},
			expected{http.StatusOK, nil, nil},
		},
		{
			"empty method",
			args{"", "", nil},
			expected{http.StatusOK, nil, nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := createTestServer(t, tt.args, tt.expected)
			defer ts.Close()

			bhh := &BindmanHTTPHelper{}
			response, body, err := bhh.request(ts.URL, tt.args.method, tt.args.contentType, tt.args.payload)
			if err != tt.expected.err {
				t.Errorf("error = %v, wantErr %v", err, tt.expected.err)
			}
			if response.StatusCode != tt.expected.code {
				t.Errorf("status code = %v, want %v", response.StatusCode, tt.expected.code)
			}
			if string(body) != string(tt.expected.body) {
				t.Errorf("body = %s, want %s", string(body), string(tt.expected.body))
			}
		})
	}

	t.Run("request time out", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second * 11)
		}))
		defer ts.Close()

		bhh := &BindmanHTTPHelper{}
		resp, data, err := bhh.request(ts.URL, "POST", "", nil)

		if resp != nil {
			t.Error("the response must be nil")
		}
		if data != nil {
			t.Error("the data value must be nil")
		}
		if !strings.Contains(err.Error(), "net/http: request canceled (Client.Timeout exceeded while awaiting headers)") {
			t.Error("the error message must contains substring 'net/http: request canceled (Client.Timeout exceeded while awaiting headers)")
		}

	})

	urlTestCases := []struct {
		name                 string
		url                  string
		method               string
		expectedErrorMessage string
	}{
		{"NewRequest empty url", "", "GET", "unsupported protocol scheme"},
		{"NewRequest unsupported protocol scheme", "unsupported://localhost", "GET", "unsupported protocol scheme"},
		{"NewRequest unsupported protocol scheme", "", "bad method", "invalid method"},
	}

	for _, tt := range urlTestCases {
		t.Run(tt.name, func(t *testing.T) {
			bhh := &BindmanHTTPHelper{}
			resp, data, err := bhh.request(tt.url, tt.method, "", nil)
			if resp != nil {
				t.Error("the response must be nil")
			}
			if data != nil {
				t.Error("the data value must be nil")
			}
			if !strings.Contains(err.Error(), tt.expectedErrorMessage) {
				t.Errorf("the error message must contains substring '%s'; got %s", tt.expectedErrorMessage, err.Error())
			}
		})
	}
}

func createTestServer(t *testing.T, args args, expected expected) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//a GET method is used when an empty string is passed by method value
		if args.method == "" {
			if r.Method != "GET" {
				t.Errorf("expected method 'GET', got %s", r.Method)
			}
		} else if r.Method != args.method {
			t.Errorf("expected method %s, got %s", args.method, r.Method)
		}
		if r.Header.Get("Content-Type") != args.contentType {
			t.Errorf("expected Content-Type header value %s, got %s", args.contentType, r.Header.Get("Content-Type"))
		}

		body, _ := ioutil.ReadAll(r.Body)

		if string(body) != string(args.payload) {
			t.Errorf("request body = %v, want %v", string(body), string(args.payload))
		}

		_, err := fmt.Fprint(w, string(expected.body))
		if err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(expected.code)
	}))
}
