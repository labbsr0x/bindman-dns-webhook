package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"

	"github.com/labbsr0x/sandman-dns-webhook/src/types"
)

// DNSWebhookClient defines the basic structure of a DNS Listener
type DNSWebhookClient struct {

	// ReverseProxyAddress the ip address of the reverse proxy that will handle the DNS redirections
	ReverseProxyAddress string

	// ManagerAddress the address of the dns manager instance
	ManagerAddress string

	// Tags a slice of strings denoting which dns records this dns agent can handle
	Tags []string
}

// New builds the client to communicate with the dns manager
func New() (*DNSWebhookClient, error) {
	rpa, err := getAddress("SANDMAN_REVERSE_PROXY_ADDRESS")
	if err != nil {
		return nil, err
	}

	ma, err := getAddress("SANDMAN_DNS_MANAGER_ADDRESS")
	if err != nil {
		return nil, err
	}

	tags := strings.Split(os.Getenv("SANDMAN_DNS_TAGS"), ",")
	return &DNSWebhookClient{
		ReverseProxyAddress: rpa,
		ManagerAddress:      ma,
		Tags:                tags,
	}, nil
}

// AddRecord is a function that calls the defined webhook to add a new dns record
func (l *DNSWebhookClient) AddRecord(name string, tags []string, ttl int) (result bool, err error) {
	record := types.DNSRecord{IPAddr: l.ReverseProxyAddress, Name: name, Tags: tags, TTL: ttl}
	ok, errs := l.checkRecord(&record)
	if ok {
		record, _ := json.Marshal(record)
		_, resp, err := PostHTTP(getRecordAPI(l.ManagerAddress, ""), record)
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
	_, resp, err := DeleteHTTP(getRecordAPI(l.ManagerAddress, name))

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

// checkRecord verifies if the tags attached to a record allows it to be handled by the registered DNSManager provider
func (l *DNSWebhookClient) checkRecord(record *types.DNSRecord) (bool, []string) {
	noErrors := false
	var errs []string

	// dumb implementation but linear O(n + m)
	rm := make(map[string]bool)
	for ri := 0; ri < len(l.Tags); ri++ {
		rm[l.Tags[ri]] = true
	}

	errs = append(errs, "No matching tags found")
	for ri := 0; ri < len(record.Tags); ri++ {
		if rm[record.Tags[ri]] {
			noErrors = true
			errs = nil
		}
	}

	ok, rErrors := record.Check()
	noErrors = noErrors && ok
	errs = append(errs, rErrors...)

	return noErrors, errs
}
