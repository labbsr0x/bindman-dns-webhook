package hook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labbsr0x/bindman-dns-webhook/src/hook/metrics"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// DNSWebhook defines the basic structure of a DNS Webhook
type DNSWebhook struct {

	// DNSManager defines the dnsmanager object this webhook will call
	DNSManager types.DNSManager
}

// Initialize starts up a dns manager webhook
func Initialize(manager types.DNSManager, serviceName, serviceVersion string) {
	hook := &DNSWebhook{manager}

	prometheus := metrics.New(serviceName, serviceVersion)

	router := mux.NewRouter()
	router.HandleFunc("/records", hook.GetDNSRecords).Methods("GET")
	router.HandleFunc("/records/{name}/{type}", hook.GetDNSRecord).Methods("GET")
	router.HandleFunc("/records/{name}/{type}", hook.RemoveDNSRecord).Methods("DELETE")
	router.HandleFunc("/records", hook.AddDNSRecord).Methods("POST")
	router.HandleFunc("/records", hook.UpdateDNSRecord).Methods("PUT")

	// exposes /metrics endpoint with standard golang metrics used by prometheus
	router.Handle("/metrics", promhttp.Handler())
	router.Use(prometheus.MetricsMiddleware)

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

// GetDNSRecord gets a specific DNS Record. DNS Record name and type comes from url params
func (m *DNSWebhook) GetDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	logrus.Infof("GetDNSRecord call. Http Request: %v", r)

	vars := mux.Vars(r)

	resp, err := m.DNSManager.GetDNSRecord(vars["name"], vars["type"])
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

	resp, err := m.DNSManager.RemoveDNSRecord(vars["name"], vars["type"])
	types.PanicIfError(types.Error{Message: fmt.Sprintf("Not possible to remove the DNS record '%s'", vars["name"]), Code: 500, Err: err})

	write200Response(resp, w)
}

// AddDNSRecord handles a POST request
// Expects a DNSRecord object as a body payload
func (m *DNSWebhook) AddDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	logrus.Infof("AddDNSRecord call. Http Request: %v", r)
	code, err := m.addOrUpdateDNSRecord(w, r, m.DNSManager.AddDNSRecord)
	types.PanicIfError(types.Error{Message: "Not possible to add a new DNS Record", Code: code, Err: err})
}

// UpdateDNSRecord updates a dns record
// Expects a DNSRecord object as a body payload
func (m *DNSWebhook) UpdateDNSRecord(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	logrus.Infof("UpdateDNSRecord call. Http Request: %v", r)
	code, err := m.addOrUpdateDNSRecord(w, r, m.DNSManager.UpdateDNSRecord)
	types.PanicIfError(types.Error{Message: "Not possible to update the DNS Record", Code: code, Err: err})
}

// actOrUpdateDNSRecord
func (m *DNSWebhook) addOrUpdateDNSRecord(w http.ResponseWriter, r *http.Request, do func(record types.DNSRecord) (bool, error)) (int, error) {
	var record types.DNSRecord
	var resp bool
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&record)
	if err == nil {
		resp, err = do(record) // call to BL provider
		if err == nil {
			write200Response(resp, w)
			return 200, nil
		}
		return 500, err
	}
	return 400, err
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
