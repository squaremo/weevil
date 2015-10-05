package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", homePage)
	router.HandleFunc("/index.html", homePage)
	router.PathPrefix("/res/").HandlerFunc(handleResource)

	router.HandleFunc("/api/{service}/", listInstances)
	router.HandleFunc("/api/", listServices)

	http.ListenAndServe("0.0.0.0:7070", router)
}

func listServices(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{})
}

func listInstances(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{})
}

func handleResource(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[1:]
	http.ServeFile(w, r, file)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
