package main

import (
	"os"
	"strings"

	"github.com/labbsr0x/sandman-dns-webhook/src/hook"
	"github.com/labbsr0x/sandman-dns-webhook/src/types"
)

func main() {
	config := GetConfig()
	manager := Bind9Manager{DNSRecords: make(map[string]types.DNSRecord)}
	hook.Initialize(config.Tags, &manager)
}

// Config holds the config information for this manager
type Config struct {
	Tags []string
}

// GetConfig reads the expected os environment variables
func GetConfig() *Config {
	return &Config{
		Tags: strings.Split(os.Getenv("TAGS"), ","),
	}
}

// Bind9Manager holds the information for managing a bing9 dns server
type Bind9Manager struct {
	DNSRecords map[string]types.DNSRecord
}

// GetDNSRecords retrieves all the dns records being managed
func (m *Bind9Manager) GetDNSRecords() ([]types.DNSRecord, error) {
	toReturn := []types.DNSRecord{}
	for _, v := range m.DNSRecords {
		toReturn = append(toReturn, v)
	}
	return toReturn, nil
}

// GetDNSRecord retrieves the dns record identified by name
func (m *Bind9Manager) GetDNSRecord(name string) (types.DNSRecord, error) {
	return m.DNSRecords[name], nil
}

// AddDNSRecord adds a new DNS record
func (m *Bind9Manager) AddDNSRecord(record types.DNSRecord) (bool, error) {
	m.DNSRecords[record.Name] = record
	return true, nil
}

// RemoveDNSRecord removes a DNS record
func (m *Bind9Manager) RemoveDNSRecord(name string) (bool, error) {
	delete(m.DNSRecords, name)
	return true, nil
}
