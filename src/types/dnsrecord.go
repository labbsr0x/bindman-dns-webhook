package types

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// DNSRecord defines what we understand as a DNSRecord
type DNSRecord struct {
	// Name the DNS host name
	Name string `json:"name"`

	// Value the value of this record
	Value string `json:"value"`

	// Type the record type
	Type string `json:"type"`
}

// Check verifies if the DNS record satisfies certain conditions
func (record *DNSRecord) Check() (bool, []string) {
	logrus.Infof("Record to check: '%v'", record)
	noErrors := true
	var errs []string

	if strings.Trim(record.Name, " ") == "" {
		noErrors = false
		errs = append(errs, "Empty record name")
	}

	if strings.Trim(record.Value, " ") == "" {
		noErrors = false
		errs = append(errs, "Empty value")
	}

	if strings.Trim(record.Type, " ") == "" {
		noErrors = false
		errs = append(errs, "Empty type")
	}

	return noErrors, errs
}
