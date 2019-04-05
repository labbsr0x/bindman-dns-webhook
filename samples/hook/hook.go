package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/labbsr0x/bindman-dns-webhook/src/hook"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
)

func main() {
	manager := DummyManager{DNSRecords: make(map[string]types.DNSRecord), TTL: 3600}

	// get ttl from env
	if ttl, err := strconv.Atoi(strings.Trim(os.Getenv("BINDMAN_DNS_TTL"), " ")); err == nil {
		manager.TTL = ttl
	}
	hook.Initialize(&manager, "1")
}

// DummyManager holds the information for managing a dummy dns server
type DummyManager struct {
	DNSRecords map[string]types.DNSRecord
	TTL        int
}

// GetDNSRecords retrieves all the dns records being managed
func (m *DummyManager) GetDNSRecords() ([]types.DNSRecord, error) {
	toReturn := []types.DNSRecord{}
	for _, v := range m.DNSRecords {
		toReturn = append(toReturn, v)
	}
	return toReturn, nil
}

// GetDNSRecord retrieves the dns record identified by name
func (m *DummyManager) GetDNSRecord(name, recordType string) (*types.DNSRecord, error) {
	if record, ok := m.DNSRecords[name]; ok {
		return &record, nil
	}
	return nil, nil
}

// AddDNSRecord adds a new DNS record
func (m *DummyManager) AddDNSRecord(record types.DNSRecord) (bool, error) {
	return m.UpdateDNSRecord(record)
}

// RemoveDNSRecord removes a DNS record
func (m *DummyManager) RemoveDNSRecord(name, recordType string) (bool, error) {
	delete(m.DNSRecords, name)
	return true, nil
}

// UpdateDNSRecord updates a DNS record. Adds, if record does not exist
func (m *DummyManager) UpdateDNSRecord(record types.DNSRecord) (bool, error) {
	m.DNSRecords[record.Name] = record
	return true, nil
}
