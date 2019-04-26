package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	hookClient "github.com/labbsr0x/bindman-dns-webhook/src/client"
	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
	"net/http"
	"os"
)

var webHookClient *hookClient.DNSWebhookClient

func main() {
	managerAddress := os.Getenv("BINDMAN_DNS_MANAGER_ADDRESS")
	client, err := hookClient.New(managerAddress, http.DefaultClient)
	if err != nil {
		panic(err)
	}
	webHookClient = client

	router := mux.NewRouter()
	router.HandleFunc("/add/{name}", addRecord).Methods("GET", "POST")
	router.HandleFunc("/remove/{name}/{type}", removeRecord).Methods("GET", "POST", "DELETE")
	router.HandleFunc("/update/{name}", updateRecord).Methods("PUT")

	http.ListenAndServe("0.0.0.0:7071", router)
}

func addRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := webHookClient.AddRecord(vars["name"], "A", "0.0.0.0")
	if err == nil {
		records, err := webHookClient.GetRecords()
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(records)
			return
		}
	}
	http.Error(w, err.Error(), 500)
}

func updateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := webHookClient.UpdateRecord(&hookTypes.DNSRecord{Name: vars["name"], Type: "A", Value: "0.0.0.0"})
	if err == nil {
		records, err := webHookClient.GetRecords()
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(records)
			return
		}
	}
	http.Error(w, err.Error(), 500)
}

func removeRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := webHookClient.RemoveRecord(vars["name"], vars["type"])
	if err == nil {
		records, err := webHookClient.GetRecords()
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(records)
			return
		}
	}
	http.Error(w, err.Error(), 500)
}
