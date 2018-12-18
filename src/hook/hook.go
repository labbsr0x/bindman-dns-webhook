package hook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"

	"github.com/gorilla/mux"
)

// DNSWebhook defines the basic structure of a DNS Manager
type DNSWebhook struct {

	// DNSManager defines the dnsmanager object this webhook will call
	DNSManager types.DNSManager
}

// Initialize starts up a dns manager webhook
func Initialize(manager types.DNSManager) {
	hook := &DNSWebhook{manager}
	router := mux.NewRouter()
	router.HandleFunc("/records", hook.GetDNSRecords).Methods("GET")
	router.HandleFunc("/records/{name}", hook.GetDNSRecord).Methods("GET")
	router.HandleFunc("/records/{name}", hook.RemoveDNSRecord).Methods("DELETE")
	router.HandleFunc("/records", hook.AddDNSRecord).Methods("POST")
	router.HandleFunc("/records", hook.UpdateDNSRecord).Methods("PUT")

	logrus.Info("Initialized DNS Manager Webhook")
	err := http.ListenAndServe("0.0.0.0:7070", router)
	if err != nil {
		logrus.Errorf("Error initializing the DNS Manager Webhook: %v", err)
	}
}

// GetDNSRecords lists the registered DNS Records
func (m *DNSWebhook) GetDNSRecords(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	logrus.Infof("GetDNSRecords call. Http Request: %v", r)

	resp, err := m.DNSManager.GetDNSRecords()
	types.PanicIfError(types.Error{Message: "Not possible to get the DNS Records", Code: 500, Err: err})

	write200Response(resp, w)
}

// GetDNSRecord gets a specific DNS Record
func (m *DNSWebhook) GetDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	logrus.Infof("GetDNSRecord call. Http Request: %v", r)

	vars := mux.Vars(r)

	resp, err := m.DNSManager.GetDNSRecord(vars["name"])
	types.PanicIfError(types.Error{Message: fmt.Sprintf("Not possible to get the DNS Record '%s'", vars["name"]), Code: 500, Err: err})

	if resp == nil {
		types.Panic(types.Error{Message: fmt.Sprintf("No record with name '%s'", vars["name"]), Code: 404, Err: nil})
	}

	write200Response(resp, w)
}

// RemoveDNSRecord removes a dns record identified by its name
func (m *DNSWebhook) RemoveDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	logrus.Infof("RemoveDNSRecord call. Http Request: %v", r)
	vars := mux.Vars(r)

	resp, err := m.DNSManager.RemoveDNSRecord(vars["name"])
	types.PanicIfError(types.Error{Message: fmt.Sprintf("Not possible to remove the DNS record '%s'", vars["name"]), Code: 500, Err: err})

	write200Response(resp, w)
}

// AddDNSRecord handles a POST request
// Expects a DNSRecord object as a body payload
func (m *DNSWebhook) AddDNSRecord(w http.ResponseWriter, r *http.Request) {
	m.addOrUpdateDNSRecord(w, r, m.DNSManager.AddDNSRecord)
}

// UpdateDNSRecord updates a dns record
// Expects a DNSRecord object as a body payload
func (m *DNSWebhook) UpdateDNSRecord(w http.ResponseWriter, r *http.Request) {
	m.addOrUpdateDNSRecord(w, r, m.DNSManager.UpdateDNSRecord)
}

// actOrUpdateDNSRecord
func (m *DNSWebhook) addOrUpdateDNSRecord(w http.ResponseWriter, r *http.Request, action func(record types.DNSRecord) (bool, error)) {
	defer handleError(w)
	logrus.Infof("UpdateDNSRecord call. Http Request: %v", r)

	decoder := json.NewDecoder(r.Body)
	var record types.DNSRecord
	err := decoder.Decode(&record)
	types.PanicIfError(types.Error{Message: fmt.Sprintf("Not possible to parse the AddDNSRecord body payload (%s)", r.Body), Code: 400, Err: err})

	resp, err := action(record) // call to BL provider
	types.PanicIfError(types.Error{Message: "Not possible to add a new DNS record", Code: 500, Err: err})

	write200Response(resp, w)
}

// write200Response writes the response to be sent
func write200Response(payload interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	err := json.NewEncoder(w).Encode(payload)
	types.PanicIfError(types.Error{Message: "Not possible to write 200 response", Code: 500, Err: err})

	logrus.Infof("200 Response sent. Payload: %s", payload)
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
