package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/labbsr0x/sandman-dns-webhook/src/client"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/add/{name}", addRecord).Methods("GET", "POST")
	router.HandleFunc("/remove/{name}", removeRecord).Methods("GET", "POST", "DELETE")

	http.ListenAndServe("0.0.0.0:7071", router)
}

func addRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	client := GetClient()
	result, err := client.AddRecord(vars["name"], []string{"test2"}, 12345)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
	}
}

func removeRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	client := GetClient()
	result, err := client.RemoveRecord(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
	}
}

// GetClient builds the client to communicate with the dns manager
func GetClient() *client.DNSWebhookClient {
	return &client.DNSWebhookClient{
		ReverseProxyAddress: os.Getenv("REVERSE_PROXY_ADDRESS"),
		ManagerAddress:      os.Getenv("MANAGER_ADDRESS"),
	}
}
