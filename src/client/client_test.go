package client

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/labbsr0x/bindman-dns-webhook/src/types"
)

func initClient() (*DNSWebhookClient, *MockHTTPHelper) {
	env := "BINDMAN_DNS_MANAGER_ADDRESS"
	mockHelper := new(MockHTTPHelper)
	os.Setenv(env, "0.0.0.0")
	c, _ := New(mockHelper)

	return c, mockHelper
}

func TestNew(t *testing.T) {
	env := "BINDMAN_DNS_MANAGER_ADDRESS"
	mockHelper := new(MockHTTPHelper)
	os.Setenv(env, "0.0.0.0")
	_, err := New(mockHelper)
	if err != nil {
		t.Errorf("Expecting client.New to succeed. Got error instead: '%v'", err)
	}

	c, err := New(nil)
	if err == nil {
		t.Errorf("Expecting client.New to fail. Got success instead: '%v'", c)
	}

	os.Setenv(env, "")
	c, err = New(mockHelper)
	if err == nil {
		t.Errorf("Expecting client.New to fail. Got success instead: '%v'", c)
	}
}

func TestGetRecords(t *testing.T) {
	c, mockHelper := initClient()
	expected := []types.DNSRecord{{}, {}}
	mockHelper.GetData, _ = json.Marshal(expected)

	records, err := c.GetRecords()
	if err != nil {
		t.Errorf("Expecting successfull execution of GetRecords. Got error instead: '%v'", err)
	}
	if len(records) != len(expected) {
		t.Errorf("Expecting the number of records to be exactly the ")
	}
}

func TestGetRecord(t *testing.T) {
	c, mockHelper := initClient()
	expected := types.DNSRecord{Name: "teste"}
	mockHelper.GetData, _ = json.Marshal(expected)

	record, err := c.GetRecord(expected.Name)
	if err != nil {
		t.Errorf("Expecting successfull execution of GetRecord. Got error instead: '%v'", err)
	}

	if record.Name != expected.Name {
		t.Errorf("Expecting the recovered record name to match exactly the expected record. Got '%v' instead", record.Name)
	}
}

func TestAddRecord(t *testing.T) {
	c, mockHelper := initClient()
	expectedRecord := types.DNSRecord{Name: "teste", Value: "0.0.0.0", Type: "A"}
	expetectedResult := true

	mockHelper.PostData, _ = json.Marshal(expetectedResult)
	result, err := c.AddRecord(expectedRecord.Name, expectedRecord.Type, expectedRecord.Value)
	if err != nil {
		t.Errorf("Expecting to successfully add the record. Got error instead: %v", err)
	}

	if result != expetectedResult {
		t.Errorf("Expecting to successfully add the record. Got failure instead.")
	}
}

func TestRemoveRecord(t *testing.T) {
	c, mockHelper := initClient()
	expetectedResult := true

	mockHelper.DeleteData, _ = json.Marshal(expetectedResult)
	result, err := c.RemoveRecord("teste")
	if err != nil {
		t.Errorf("Expecting to successfully add the record. Got error instead: %v", err)
	}

	if result != expetectedResult {
		t.Errorf("Expecting to successfully add the record. Got failure instead.")
	}
}

func TestGetRecordAPI(t *testing.T) {
	api := getRecordAPI("manager.test.com", "test.test.com")
	expected := "http://manager.test.com/records/test.test.com"
	if api != expected {
		t.Errorf("Expecting '%v'; Got '%v'", expected, api)
	}
}

func TestGetAddress(t *testing.T) {
	env := "BINDMAN_TEST_ADDRESS"
	os.Setenv(env, "test.com")

	addr, err := getAddress(env)
	if addr == "" || err != nil {
		t.Errorf("Expecting the getAddress func to succeed. Got err instead: '%v'", err)
	}

	os.Setenv(env, "http://test.com")
	addr, err = getAddress(env)
	if err == nil {
		t.Errorf("Expecting the getAddress func to return error. Got success instead: '%v'", addr)
	}
}

type MockHTTPHelper struct {
	PutData    []byte
	PostData   []byte
	GetData    []byte
	DeleteData []byte
}

func (m *MockHTTPHelper) Put(url string, data []byte) (*http.Response, []byte, error) {
	return &http.Response{}, m.PutData, nil
}

func (m *MockHTTPHelper) Post(url string, data []byte) (*http.Response, []byte, error) {
	return &http.Response{}, m.PostData, nil
}

func (m *MockHTTPHelper) Get(url string) (*http.Response, []byte, error) {
	return &http.Response{}, m.GetData, nil
}
func (m *MockHTTPHelper) Delete(url string) (*http.Response, []byte, error) {
	return &http.Response{}, m.DeleteData, nil
}
