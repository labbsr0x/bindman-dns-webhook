package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/labbsr0x/bindman-dns-webhook/src/types"
)

// DNSWebhookClient defines the basic structure of a DNS Listener
type DNSWebhookClient struct {

	// ManagerAddress the address of the dns manager instance
	ManagerAddress string

	http HTTPHelper
}

// New builds the client to communicate with the dns manager
func New(httpHelper HTTPHelper) (*DNSWebhookClient, error) {
	ma, err := getAddress("BINDMAN_DNS_MANAGER_ADDRESS")
	if err != nil {
		return nil, err
	}

	if httpHelper == nil {
		return nil, fmt.Errorf("Not possible to start a listener without an HTTPHelper instance")
	}

	return &DNSWebhookClient{
		ManagerAddress: ma,
		http:           httpHelper,
	}, nil
}

// AddRecord adds a DNS record
func (l *DNSWebhookClient) AddRecord(name string, recordType string, value string) (result bool, err error) {
	record := &types.DNSRecord{Value: value, Name: name, Type: recordType}
	ok, errs := record.Check()
	if ok {
		record, _ := json.Marshal(record)
		_, resp, err := l.http.Post(getRecordAPI(l.ManagerAddress, ""), record)
		if err != nil {
			logrus.Errorf("ERR: %s", err)
			return false, err
		}

		err = json.Unmarshal(resp, &result)
		return result, err
	}
	return false, fmt.Errorf("Invalid DNS Record: %v", strings.Join(errs, ", "))
}

// RemoveRecord is a function that calls the defined webhook to remove a specific dns record
func (l *DNSWebhookClient) RemoveRecord(name string) (result bool, err error) {
	_, resp, err := l.http.Delete(getRecordAPI(l.ManagerAddress, name))

	if err != nil {
		logrus.Errorf("ERR: %s", err)
		return false, err
	}

	json.Unmarshal(resp, &result)
	return result, err
}

// getRecordAPI builds the url for consuming the api
func getRecordAPI(managerAddress string, params string) string {
	u, _ := url.Parse("http://" + managerAddress)
	u.Path = path.Join(u.Path, "/records/", params)
	return u.String()
}

// getAddress gets an env variable address identified by name
func getAddress(name string) (addr string, err error) {
	addr = os.Getenv(name)
	addr = strings.Trim(addr, " ")

	if addr == "" {
		err = fmt.Errorf("The %s environment variable was not defined", name)
	}

	if strings.Contains(addr, "http") {
		err = fmt.Errorf("The %s environment variable cannot have a schema defined", name)
	}

	return
}
