package types

// DNSRecord defines what we understand as a DNSRecord
type DNSRecord struct {
	// Name the DNS host name
	Name string `json:"name"`

	// IPAddr the ip address of the host (usually the load balancers ip)
	IPAddr string `json:"ipaddr"`

	// TTL the time to live of this dns record
	TTL int `json:"ttl"`

	// Tags slice of strings that identifies this dns record
	Tags []string `json:"tags"`
}
