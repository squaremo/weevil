package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{service}/", listInstances)
	router.HandleFunc("/", listServices)
	http.ListenAndServe("0.0.0.0:7070", router)
}

func listServices(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{})
}

func listInstances(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{})
}
