package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	hookClient "github.com/labbsr0x/bindman-dns-webhook/src/client"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/add/{name}", addRecord).Methods("GET", "POST")
	router.HandleFunc("/remove/{name}", removeRecord).Methods("GET", "POST", "DELETE")

	http.ListenAndServe("0.0.0.0:7071", router)
}

func addRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	client, err := hookClient.New()
	if err == nil {
		result, err := client.AddRecord(vars["name"], "A", "0.0.0.0")
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(result)
			return
		}
	}
	http.Error(w, err.Error(), 500)
}

func removeRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	client, err := hookClient.New()
	if err == nil {
		result, err := client.RemoveRecord(vars["name"])
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(result)
			return
		}
	}
	http.Error(w, err.Error(), 500)
}
