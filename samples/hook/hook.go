package main

import (
	"github.com/labbsr0x/sandman-dns-webhook/src/hook"
	"github.com/labbsr0x/sandman-dns-webhook/src/types"
)

func main() {
	manager := DummyManager{DNSRecords: make(map[string]types.DNSRecord)}
	hook.Initialize(&manager)
}

// DummyManager holds the information for managing a dummy dns server
type DummyManager struct {
	DNSRecords map[string]types.DNSRecord
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
func (m *DummyManager) GetDNSRecord(name string) (*types.DNSRecord, error) {
	if record, ok := m.DNSRecords[name]; ok {
		return &record, nil
	}
	return nil, nil
}

// AddDNSRecord adds a new DNS record
func (m *DummyManager) AddDNSRecord(record types.DNSRecord) (bool, error) {
	m.DNSRecords[record.Name] = record
	return true, nil
}

// RemoveDNSRecord removes a DNS record
func (m *DummyManager) RemoveDNSRecord(name string) (bool, error) {
	delete(m.DNSRecords, name)
	return true, nil
}
