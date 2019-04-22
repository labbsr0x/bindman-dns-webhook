package client

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/labbsr0x/goh/gohclient"
)

func TestNew(t *testing.T) {
	type args struct {
		managerAddress string
		http           *http.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty manager address string",
			args:    args{managerAddress: ""},
			wantErr: true,
		},
		{
			name:    "valid manager address",
			args:    args{managerAddress: "0.0.0.0"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.managerAddress, tt.args.http)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("New() must return a nil *DNSWebhookClient when an error occurred, got %v", got)
			}
		})
	}
}

func TestDNSWebhookClient_GetRecords(t *testing.T) {
	expected := []types.DNSRecord{{}, {}}
	expectedData, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := types.BadRequestError("get records error", nil)
	expectedErrorData, err := json.Marshal(expectedError)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		clientAPI gohclient.API
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult interface{}
		wantErr    bool
	}{
		{
			name:       "request success and 200 status code",
			fields:     fields{&MockHTTPHelperSuccess{Status: http.StatusOK, Data: expectedData}},
			wantResult: expected,
			wantErr:    false,
		},
		{
			name:       "request success and 400 status code",
			fields:     fields{&MockHTTPHelperSuccess{Status: http.StatusBadRequest, Data: expectedErrorData}},
			wantResult: expectedError,
			wantErr:    true,
		},
		{
			name:       "request error",
			fields:     fields{&MockHTTPHelperError{err: &url.Error{Op: "request error get record"}}},
			wantResult: &url.Error{Op: "request error get record"},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &DNSWebhookClient{
				ClientAPI: tt.fields.clientAPI,
			}
			gotResult, err := l.GetRecords()
			if (err != nil) != tt.wantErr {
				t.Errorf("DNSWebhookClient.GetRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// compare error with wanted result when error expected
			if tt.wantErr {
				if !reflect.DeepEqual(err, tt.wantResult) {
					t.Errorf("DNSWebhookClient.GetRecords() = %#v, want %#v", err, tt.wantResult)
				}
			} else {
				if !reflect.DeepEqual(gotResult, tt.wantResult) {
					t.Errorf("DNSWebhookClient.GetRecords() = %v, want %v", gotResult, tt.wantResult)
				}
			}
		})
	}
}

func TestDNSWebhookClient_GetRecord(t *testing.T) {
	expected := types.DNSRecord{Name: "test", Type: "A"}
	expectedData, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := types.BadRequestError("get record error", nil)
	expectedErrorData, err := json.Marshal(expectedError)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		clientAPI gohclient.API
	}
	type parans struct {
		name string
		typ  string
	}
	tests := []struct {
		name       string
		fields     fields
		parans     parans
		wantResult interface{}
		wantErr    bool
	}{
		{
			name:       "request success and 200 status code",
			fields:     fields{&MockHTTPHelperSuccess{Status: http.StatusOK, Data: expectedData}},
			parans:     parans{name: "teste", typ: "A"},
			wantResult: expected,
			wantErr:    false,
		},
		{
			name:       "request success and 400 status code",
			fields:     fields{&MockHTTPHelperSuccess{Status: http.StatusBadRequest, Data: expectedErrorData}},
			wantResult: expectedError,
			wantErr:    true,
		},
		{
			name:       "request error",
			fields:     fields{&MockHTTPHelperError{err: &url.Error{Op: "request error get record"}}},
			wantResult: &url.Error{Op: "request error get record"},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &DNSWebhookClient{
				ClientAPI: tt.fields.clientAPI,
			}
			gotResult, err := l.GetRecord(tt.parans.name, tt.parans.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("DNSWebhookClient.GetRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// compare error with wanted result when error expected
			if tt.wantErr {
				if !reflect.DeepEqual(err, tt.wantResult) {
					t.Errorf("DNSWebhookClient.GetRecord() = %#v, want %#v", err, tt.wantResult)
				}
			} else {
				if !reflect.DeepEqual(gotResult, tt.wantResult) {
					t.Errorf("DNSWebhookClient.GetRecord() = %v, want %v", gotResult, tt.wantResult)
				}
			}
		})
	}
}

func TestDNSWebhookClient_RemoveRecord(t *testing.T) {
	expectedError := types.BadRequestError("remove record error", nil)
	expectedErrorData, err := json.Marshal(expectedError)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		clientAPI gohclient.API
	}
	type args struct {
		name       string
		recordType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name:   "request success and 204 status code",
			fields: fields{&MockHTTPHelperSuccess{Status: http.StatusNoContent}},
			args:   args{name: "teste", recordType: "A"},
		},
		{
			name:    "request success and 400 status code",
			fields:  fields{&MockHTTPHelperSuccess{Status: http.StatusBadRequest, Data: expectedErrorData}},
			wantErr: expectedError,
		},
		{
			name:    "request error",
			fields:  fields{&MockHTTPHelperError{err: &url.Error{Op: "request error remove records"}}},
			wantErr: &url.Error{Op: "request error remove records"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &DNSWebhookClient{
				ClientAPI: tt.fields.clientAPI,
			}
			if err := l.RemoveRecord(tt.args.name, tt.args.recordType); (err != nil) != (tt.wantErr != nil) {
				t.Errorf("DNSWebhookClient.RemoveRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDNSWebhookClient_UpdateRecord(t *testing.T) {
	expectedError := types.BadRequestError("add record error", nil)
	expectedErrorData, err := json.Marshal(expectedError)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		clientAPI gohclient.API
	}
	type args struct {
		name       string
		recordType string
		value      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "request success and 204 status code",
			fields: fields{&MockHTTPHelperSuccess{Status: http.StatusNoContent}},
			args:   args{name: "test", recordType: "A", value: "A"},
		},
		{
			name:    "error check record - do not execute request",
			fields:  fields{&MockHTTPHelperSuccess{}},
			args:    args{recordType: "A", value: "A"},
			wantErr: true,
		},
		{
			name:    "request success and 400 status code",
			fields:  fields{&MockHTTPHelperSuccess{Status: http.StatusBadRequest, Data: expectedErrorData}},
			wantErr: true,
		},
		{
			name:    "request error",
			fields:  fields{&MockHTTPHelperError{err: &url.Error{Op: "request error add record"}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &DNSWebhookClient{
				ClientAPI: tt.fields.clientAPI,
			}
			if err := l.UpdateRecord(&types.DNSRecord{Name: tt.args.name, Type: tt.args.recordType, Value: tt.args.value}); (err != nil) != tt.wantErr {
				t.Errorf("DNSWebhookClient.AddRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDNSWebhookClient_AddRecord(t *testing.T) {
	expectedError := types.BadRequestError("add record error", nil)
	expectedErrorData, err := json.Marshal(expectedError)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		clientAPI gohclient.API
	}
	type args struct {
		name       string
		recordType string
		value      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "request success and 204 status code",
			fields: fields{&MockHTTPHelperSuccess{Status: http.StatusNoContent}},
			args:   args{name: "test", recordType: "A", value: "A"},
		},
		{
			name:    "error check record - do not execute request",
			fields:  fields{&MockHTTPHelperSuccess{}},
			args:    args{recordType: "A", value: "A"},
			wantErr: true,
		},
		{
			name:    "request success and 400 status code",
			fields:  fields{&MockHTTPHelperSuccess{Status: http.StatusBadRequest, Data: expectedErrorData}},
			args:    args{name: "test", recordType: "A", value: "A"},
			wantErr: true,
		},
		{
			name:    "request error",
			fields:  fields{&MockHTTPHelperError{err: &url.Error{Op: "request error add record"}}},
			args:    args{name: "test", recordType: "A", value: "A"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &DNSWebhookClient{
				ClientAPI: tt.fields.clientAPI,
			}
			if err := l.AddRecord(tt.args.name, tt.args.recordType, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("DNSWebhookClient.AddRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type MockHTTPHelperSuccess struct {
	Data   []byte
	Status int
}

func (m *MockHTTPHelperSuccess) Put(url string, data []byte) (*http.Response, []byte, error) {
	return &http.Response{StatusCode: m.Status}, m.Data, nil
}

func (m *MockHTTPHelperSuccess) Post(url string, data []byte) (*http.Response, []byte, error) {
	return &http.Response{StatusCode: m.Status}, m.Data, nil
}

func (m *MockHTTPHelperSuccess) Get(url string) (*http.Response, []byte, error) {
	return &http.Response{StatusCode: m.Status}, m.Data, nil
}
func (m *MockHTTPHelperSuccess) Delete(url string) (*http.Response, []byte, error) {
	return &http.Response{StatusCode: m.Status}, m.Data, nil
}

type MockHTTPHelperError struct {
	err error
}

func (m MockHTTPHelperError) Put(url string, data []byte) (*http.Response, []byte, error) {
	return nil, nil, m.err
}

func (m MockHTTPHelperError) Post(url string, data []byte) (*http.Response, []byte, error) {
	return nil, nil, m.err
}

func (m MockHTTPHelperError) Get(url string) (*http.Response, []byte, error) {
	return nil, nil, m.err
}

func (m MockHTTPHelperError) Delete(url string) (*http.Response, []byte, error) {
	return nil, nil, m.err
}
