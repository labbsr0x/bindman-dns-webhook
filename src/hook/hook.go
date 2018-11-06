package hook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labbsr0x/sandman-dns-webhook/src/types"

	"github.com/gorilla/mux"
)

// DNSWebhook defines the basic structure of a DNS Manager
type DNSWebhook struct {

	// Tags a slice of strings denoting which dns records this webhook can handle
	Tags []string

	// DNSManager defines the dnsmanager object this webhook will call
	DNSManager types.DNSManager
}

// Initialize starts up a dns manager webhook
func Initialize(tags []string, manager types.DNSManager) error {
	hook := DNSWebhook{tags, manager}
	router := mux.NewRouter()
	router.HandleFunc("/records", hook.GetDNSRecords).Methods("GET")
	router.HandleFunc("/records/{name}", hook.GetDNSRecord).Methods("GET")
	router.HandleFunc("/records", hook.AddDNSRecord).Methods("POST")
	router.HandleFunc("/records/{name}", hook.RemoveDNSRecord).Methods("DELETE")

	err := http.ListenAndServe("0.0.0.0:7070", router)
	if err != nil {
		return err
	}
	return nil
}

// GetDNSRecords lists the registered DNS Records
func (m *DNSWebhook) GetDNSRecords(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)

	resp, err := m.DNSManager.GetDNSRecords()
	types.PanicIfError(types.Error{Message: "Not possible to get the DNS Records", Code: 500, Err: err})

	write200Response(resp, w)
}

// GetDNSRecord gets a specific DNS Record
func (m *DNSWebhook) GetDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	vars := mux.Vars(r)

	resp, err := m.DNSManager.GetDNSRecord(vars["name"])
	types.PanicIfError(types.Error{Message: fmt.Sprintf("Not possible to get the DNS Record '%s'", vars["name"]), Code: 500, Err: err})

	write200Response(resp, w)
}

// AddDNSRecord handles a POST request
// Expects a DNSRecord object as a body payload
func (m *DNSWebhook) AddDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	decoder := json.NewDecoder(r.Body)
	var record types.DNSRecord
	err := decoder.Decode(&record)
	types.PanicIfError(types.Error{Message: fmt.Sprintf("Not possible to parse the AddDNSRecord body payload (%s)", r.Body), Code: 400, Err: err})

	toReturn := false
	if m.canAddRecord(&record) {
		resp, err := m.DNSManager.AddDNSRecord(record) // call to BL provider
		types.PanicIfError(types.Error{Message: "Not possible to add a new DNS record", Code: 500, Err: err})
		toReturn = resp
	}

	write200Response(toReturn, w)
}

// RemoveDNSRecord removes a dns record identified by its name
func (m *DNSWebhook) RemoveDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	vars := mux.Vars(r)

	resp, err := m.DNSManager.RemoveDNSRecord(vars["name"])
	types.PanicIfError(types.Error{Message: fmt.Sprintf("Not possible to remove the DNS record '%s'", vars["name"]), Code: 500, Err: err})

	write200Response(resp, w)
}

// canAddRecord verifies if the tags attached to a record allows it to be handled by the registered DNSManager provider
func (m *DNSWebhook) canAddRecord(record *types.DNSRecord) bool {
	// dumb implementation but linear O(n + m)
	rm := make(map[string]bool)
	for ri := 0; ri < len(m.Tags); ri++ {
		rm[m.Tags[ri]] = true
	}

	for ri := 0; ri < len(record.Tags); ri++ {
		if rm[record.Tags[ri]] {
			return true
		}
	}

	return false
}

// write200Response writes the response to be sent
func write200Response(payload interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	err := json.NewEncoder(w).Encode(payload)
	types.PanicIfError(types.Error{Message: "Not possible to write 200 response", Code: 500, Err: err})

	logrus.Debugf("200 Response sent. Payload: %s", payload)
}

// handleError recovers from a panic
func handleError(w http.ResponseWriter) {
	r := recover()
	if r != nil {
		if err, ok := r.(types.Error); ok {
			logrus.Error(err)
			http.Error(w, err.Message, err.Code)
		} else {
			logrus.Error(r)
			http.Error(w, "Erro interno", 500)
		}
	}
}
