package hook

import (
	"fmt"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_writeJSONResponse(t *testing.T) {
	type test struct {
		Message string `json:"message"`
	}

	type args struct {
		payload    interface{}
		statusCode int
	}
	type expected struct {
		code int
		body string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{"must have empty string on body when payload = nil",
			args{nil, http.StatusOK},
			expected{http.StatusOK, ""},
		},
		{"must have Json formatted body when a non-nil payload",
			args{test{"Hello world"}, http.StatusOK},
			expected{http.StatusOK, "{\"message\":\"Hello world\"}\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			writeJSONResponse(tt.args.payload, tt.args.statusCode, res)
			if tt.expected.code != res.Code {
				t.Errorf("expected code %d, got %d", tt.expected.code, res.Code)
			}
			resBody := res.Body.String()
			if tt.expected.body != resBody {
				t.Errorf("expected code %s, got %s", tt.expected.body, resBody)
			}
			if tt.expected.body != "" {
				if res.Header().Get("Content-Type") != "application/json" {
					t.Error("the content type header value must be 'application/json' when non-empty body")
				}
			}
		})
	}
}

func Test_handleError(t *testing.T) {
	type args struct {
		e error
	}
	type expected struct {
		code int
		body string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{"do not write any response when recover function return nil",
			args{nil},
			expected{httptest.NewRecorder().Code, httptest.NewRecorder().Body.String()},
		},
		{"must write status code and body from error recovered when it is instance of types.Error ",
			args{types.BadRequestError("error", nil)},
			expected{http.StatusBadRequest, "{\"message\":\"error\",\"code\":400,\"details\":null}\n"},
		},
		{"must write a default internal server error when recovered error is not instance of types.Error ",
			args{fmt.Errorf("unknow error")},
			expected{http.StatusInternalServerError, "{\"message\":\"An internal server error occurred, please contact the system administrator.\",\"code\":500,\"details\":null}\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			defer func() {
				if tt.expected.code != res.Code {
					t.Errorf("expected code %d, got %d", tt.expected.code, res.Code)
				}
				resBody := res.Body.String()
				if tt.expected.body != resBody {
					t.Errorf("expected code %s, got %s", tt.expected.body, resBody)
				}
				if tt.expected.body != "" {
					if res.Header().Get("Content-Type") != "application/json" {
						t.Error("the content type header value must be 'application/json' when non-empty body")
					}
				}
			}()
			defer handleError(res)
			panic(tt.args.e)
		})
	}
}
