package client

import (
	"encoding/json"
	"net/url"
	"path"

	"github.com/Sirupsen/logrus"

	"github.com/labbsr0x/sandman-dns-webhook/src/types"
)

// DNSWebhookClient defines the basic structure of a DNS Listener
type DNSWebhookClient struct {

	// ReverseProxyAddress the ip address of the reverse proxy that will handle the DNS redirections
	ReverseProxyAddress string

	// WebhookAddress the address of the dns manager instance
	WebhookAddress string
}

// AddRecord is a function that calls the defined webhook to add a new dns record
func (l *DNSWebhookClient) AddRecord(name string, tags []string) (result bool, err error) {
	record, _ := json.Marshal(types.DNSRecord{IPAddr: l.ReverseProxyAddress, Name: name, Tags: tags})
	_, resp, err := PostHTTP(l.WebhookAddress, record)
	if err != nil {
		logrus.Errorf("ERR: %s", err)
		return false, err
	}

	json.Unmarshal(resp, &result)
	return result, err
}

// RemoveRecord is a function that calls the defined webhook to remove a specific dns record
func (l *DNSWebhookClient) RemoveRecord(name string) (result bool, err error) {
	u, _ := url.Parse(l.WebhookAddress)
	u.Path = path.Join(u.Path, name)
	_, resp, err := DeleteHTTP(u.String())

	if err != nil {
		logrus.Errorf("ERR: %s", err)
		return false, err
	}

	json.Unmarshal(resp, &result)
	return result, err
}
