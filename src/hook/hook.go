package hook

import (
	"encoding/json"
	"errors"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// DNSWebhook defines the basic structure of a DNS Webhook
type DNSWebhook struct {

	// DNSManager defines the dnsmanager object this webhook will call
	DNSManager types.DNSManager
}

// Initialize starts up a dns manager webhook
func Initialize(manager types.DNSManager) {
	if manager == nil {
		panic(errors.New("A non-nil DNSManager is required to initialize the hook"))
	}

	hook := &DNSWebhook{manager}
	router := mux.NewRouter()
	router.HandleFunc("/records", hook.GetDNSRecords).Methods("GET")
	router.HandleFunc("/records/{name}/{type}", hook.GetDNSRecord).Methods("GET")
	router.HandleFunc("/records/{name}/{type}", hook.RemoveDNSRecord).Methods("DELETE")
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
	types.PanicIfError(err)
	writeJSONResponse(resp, http.StatusOK, w)
}

// GetDNSRecord gets a specific DNS Record. DNS Record name and type comes from url params
func (m *DNSWebhook) GetDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	logrus.Infof("GetDNSRecord call. Http Request: %v", r)

	vars := mux.Vars(r)

	resp, err := m.DNSManager.GetDNSRecord(vars["name"], vars["type"])
	types.PanicIfError(err)
	writeJSONResponse(resp, http.StatusOK, w)
}

// RemoveDNSRecord removes a dns record identified by its name
func (m *DNSWebhook) RemoveDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	logrus.Infof("RemoveDNSRecord call. Http Request: %v", r)
	vars := mux.Vars(r)

	err := m.DNSManager.RemoveDNSRecord(vars["name"], vars["type"])
	types.PanicIfError(err)

	w.WriteHeader(http.StatusNoContent)
}

// AddDNSRecord handles a POST request
// Expects a DNSRecord object as a body payload
func (m *DNSWebhook) AddDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	logrus.Infof("AddDNSRecord call. Http Request: %v", r)
	err := m.addOrUpdateDNSRecord(w, r, m.DNSManager.AddDNSRecord)
	types.PanicIfError(err)
}

// UpdateDNSRecord updates a dns record
// Expects a DNSRecord object as a body payload
func (m *DNSWebhook) UpdateDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	logrus.Infof("UpdateDNSRecord call. Http Request: %v", r)
	err := m.addOrUpdateDNSRecord(w, r, m.DNSManager.UpdateDNSRecord)
	types.PanicIfError(err)
}

// actOrUpdateDNSRecord
func (m *DNSWebhook) addOrUpdateDNSRecord(w http.ResponseWriter, r *http.Request, do func(record types.DNSRecord) error) error {
	var record types.DNSRecord
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&record)
	if err != nil {
		if err == io.EOF {
			return types.BadRequestError("You must pass a JSON formatted record on request body", nil)
		}
		return err
	}
	errs := record.Check()
	if errs != nil {
		return types.BadRequestError("You must pass a JSON formatted record on request body", errs)
	}
	err = do(record) // call to BL provider
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
