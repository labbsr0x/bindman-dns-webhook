package types

import "strings"

// DNSRecord defines what we understand as a DNSRecord
type DNSRecord struct {
	// Name the DNS host name
	Name string `json:"name"`

	// IPAddr the ip address of the host (usually the load balancers ip)
	IPAddr string `json:"ipaddr"`
}

// Check verifies if the DNS record satisfies certain conditions
func (record *DNSRecord) Check() (bool, []string) {
	noErrors := true
	var errs []string

	if strings.Trim(record.Name, " ") == "" {
		noErrors = false
		errs = append(errs, "Empty record name")
	}

	if strings.Trim(record.IPAddr, " ") == "" {
		noErrors = false
		errs = append(errs, "Empty ip address")
	}

	return noErrors, errs
}
