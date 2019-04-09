package hook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

var records = []types.DNSRecord{{"test.com.br", "127.0.0.1", "A"}}

func TestInitialize(t *testing.T) {
	t.Run("initialize the hook with a nil DNSManager", func(t *testing.T) {
		defer func() {
			err := recover()
			if err == nil {
				t.Error("A panic must occur when a nil DNSManager is passed to Initialize function")
			}
		}()
		Initialize(nil, "1")
	})
}

func TestDNSRecordsHandlers(t *testing.T) {
	var (
		errorBadRequest       = &types.Error{Message: "test message", Code: http.StatusBadRequest}
		hookSuccess           = &DNSWebhook{&SuccessDNSManagerMock{records}}
		hookError             = &DNSWebhook{&ErrorDNSManagerMock{errorBadRequest}}
		invalidRequestBodyMsg = "Invalid request body. You must pass a JSON formatted record on request body"
	)
	type expected struct {
		code int
		body interface{}
	}

	type req struct {
		method string
		path   string
		body   interface{}
	}

	testCases := []struct {
		name     string
		req      req
		path     string
		handle   func(http.ResponseWriter, *http.Request)
		expected expected
	}{
		{"GetDNSRecords retrieving all records",
			req{},
			"",
			hookSuccess.GetDNSRecords,
			expected{http.StatusOK, records},
		},
		{"GetDNSRecords error retrieving records",
			req{},
			"",
			hookError.GetDNSRecords,
			expected{http.StatusBadRequest, errorBadRequest},
		},
		{"GetDNSRecords retrieving record",
			req{path: fmt.Sprintf("/%s/%s", records[0].Name, records[0].Type)},
			"/{name}/{type}",
			hookSuccess.GetDNSRecord,
			expected{http.StatusOK, records[0]},
		},
		{"GetDNSRecord error retrieving record",
			req{path: fmt.Sprintf("/%s/%s", records[0].Name, records[0].Type)},
			"/{name}/{type}",
			hookError.GetDNSRecord,
			expected{http.StatusBadRequest, errorBadRequest},
		},
		{"RemoveDNSRecord deleting record",
			req{path: fmt.Sprintf("/%s/%s", records[0].Name, records[0].Type)},
			"/{name}/{type}",
			hookSuccess.RemoveDNSRecord,
			expected{http.StatusNoContent, nil},
		},
		{"RemoveDNSRecord error deleting record",
			req{path: fmt.Sprintf("/%s/%s", records[0].Name, records[0].Type)},
			"/{name}/{type}",
			hookError.RemoveDNSRecord,
			expected{http.StatusBadRequest, errorBadRequest},
		},
		{"UpdateDNSRecord update record",
			req{body: records[0]},
			"",
			hookSuccess.UpdateDNSRecord,
			expected{http.StatusNoContent, nil},
		},
		{"UpdateDNSRecord error updating record",
			req{body: records[0]},
			"",
			hookError.UpdateDNSRecord,
			expected{http.StatusBadRequest, errorBadRequest},
		},
		{"UpdateDNSRecord error empty requestBody",
			req{},
			"",
			hookError.UpdateDNSRecord,
			expected{http.StatusBadRequest, types.BadRequestError(invalidRequestBodyMsg, nil)},
		},
		{"UpdateDNSRecord error invalid record on requestBody",
			req{body: types.DNSRecord{Name: "test.com.br"}},
			"",
			hookError.UpdateDNSRecord,
			expected{http.StatusBadRequest, types.BadRequestError(invalidRequestBodyMsg, nil, (&types.DNSRecord{Name: "test.com.br"}).Check()...)},
		},
		{"UpdateDNSRecord error invalid content on requestBody",
			req{body: "invalid format"},
			"",
			hookError.UpdateDNSRecord,
			expected{http.StatusBadRequest, types.BadRequestError(invalidRequestBodyMsg, nil)},
		},
		{"AddDNSRecord add record",
			req{body: records[0]},
			"",
			hookSuccess.AddDNSRecord,
			expected{http.StatusNoContent, nil},
		},
		{"AddDNSRecord error adding record",
			req{body: records[0]},
			"",
			hookError.AddDNSRecord,
			expected{http.StatusBadRequest, errorBadRequest},
		},
		{"AddDNSRecord error empty requestBody",
			req{},
			"",
			hookError.AddDNSRecord,
			expected{http.StatusBadRequest, types.BadRequestError(invalidRequestBodyMsg, nil)},
		},
		{"AddDNSRecord error invalid record on requestBody",
			req{body: types.DNSRecord{Name: "test.com.br"}},
			"",
			hookError.AddDNSRecord,
			expected{http.StatusBadRequest, types.BadRequestError(invalidRequestBodyMsg, nil, (&types.DNSRecord{Name: "test.com.br"}).Check()...)},
		},
		{"AddDNSRecord error invalid content on requestBody",
			req{body: "invalid format"},
			"",
			hookError.AddDNSRecord,
			expected{http.StatusBadRequest, types.BadRequestError(invalidRequestBodyMsg, nil)},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//prepare request body
			var buf bytes.Buffer
			if testCase.req.body != nil {
				err := json.NewEncoder(&buf).Encode(testCase.req.body)
				if err != nil {
					t.Fatal(err)
				}
			}
			//create request
			req, err := http.NewRequest(testCase.req.method, fmt.Sprintf("/records%s", testCase.req.path), &buf)
			if err != nil {
				t.Fatal(err)
			}
			//create response
			res := httptest.NewRecorder()

			// Need to create a router that we can pass the request through so that the vars will be added to the context
			router := mux.NewRouter()
			router.HandleFunc(fmt.Sprintf("/records%s", testCase.path), testCase.handle)
			//execute handler
			router.ServeHTTP(res, req)

			//create response body
			var respBody bytes.Buffer
			if testCase.expected.body != nil {
				if err := json.NewEncoder(&respBody).Encode(testCase.expected.body); err != nil {
					t.Fatal(err)
				}
				if res.Header().Get("Content-Type") != "application/json" {
					t.Errorf("handler returned unexpected content type header value: want = %s, got %s", "application/json", res.Header().Get("Content-Type"))
				}
			}

			if res.Body.String() != respBody.String() {
				t.Errorf("handler returned unexpected body: want = %s, got %s", respBody.String(), res.Body.String())
			}
			if res.Code != testCase.expected.code {
				t.Errorf("handler returned unexpected status code: want = %d, got %d", testCase.expected.code, res.Code)
			}
		})
	}

}

type SuccessDNSManagerMock struct {
	records []types.DNSRecord
}

func (m *SuccessDNSManagerMock) GetDNSRecords() ([]types.DNSRecord, error) {
	return m.records, nil
}

func (m *SuccessDNSManagerMock) GetDNSRecord(name, recordType string) (*types.DNSRecord, error) {
	record := m.records[0]
	if name == record.Name && recordType == record.Type {
		return &record, nil
	}
	return nil, types.InternalServerError(fmt.Sprintf("expected name = %s and type %s on path parameter, got name = %s and type %s", record.Name, record.Type, name, recordType), nil)
}

func (m *SuccessDNSManagerMock) RemoveDNSRecord(name, recordType string) error {
	return nil
}

func (m *SuccessDNSManagerMock) AddDNSRecord(record types.DNSRecord) error {
	return nil
}

func (m *SuccessDNSManagerMock) UpdateDNSRecord(record types.DNSRecord) error {
	return nil
}

type ErrorDNSManagerMock struct {
	error *types.Error
}

func (m *ErrorDNSManagerMock) GetDNSRecords() ([]types.DNSRecord, error) {
	return nil, m.error
}

func (m *ErrorDNSManagerMock) GetDNSRecord(name, recordType string) (*types.DNSRecord, error) {
	return nil, m.error
}

func (m *ErrorDNSManagerMock) RemoveDNSRecord(name, recordType string) error {
	return m.error
}

func (m *ErrorDNSManagerMock) AddDNSRecord(record types.DNSRecord) error {
	return m.error
}

func (m *ErrorDNSManagerMock) UpdateDNSRecord(record types.DNSRecord) error {
	return m.error
}
