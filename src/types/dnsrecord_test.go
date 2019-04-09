package types

import "testing"

func TestCheckDNSRecord(t *testing.T) {
	invalidName := "the value of field 'name' cannot be empty"
	invalidValue := "the value of field 'value' cannot be empty"
	invalidType := "the value of field 'type' cannot be empty"

	testCases := []struct {
		name     string
		record   DNSRecord
		expected []string
	}{
		{"Return nil when all attributes are ok", DNSRecord{Name: "t.test.com", Value: "0.0.0.0", Type: "A"}, nil},
		{"validate empty value", DNSRecord{Name: "", Value: "0.0.0.0", Type: "A"}, []string{invalidName}},
		{"validate empty values", DNSRecord{Name: "", Value: "", Type: "A"}, []string{invalidName, invalidValue}},
		{"validate nil attribute", DNSRecord{Name: "", Value: ""}, []string{invalidName, invalidValue, invalidType}},
		{"validate spaces values", DNSRecord{Name: " ", Value: " ", Type: " "}, []string{invalidName, invalidValue, invalidType}},
		{"validate nil attribute and empty values", DNSRecord{Value: "  ", Type: " "}, []string{invalidName, invalidValue, invalidType}},
		{"validate all nil", DNSRecord{}, []string{invalidName, invalidValue, invalidType}},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			errs := test.record.Check()
			if len(errs) != len(test.expected) {
				t.Errorf("The error array length must be %d but got %d", len(test.expected), len(errs))
				t.FailNow()
			}
			for i, err := range test.expected {
				if errs[i] != err {
					t.Errorf("Expected message was %s but got %s", err, errs[i])
				}
			}
		})
	}
}
