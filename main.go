package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func muxRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")
	// declare static file directory
	staticFileDirectory := http.Dir("./assets/")

	staticFileHandler := http.StripPrefix("/assets", http.FileServer(staticFileDirectory))

	r.PathPrefix("/assets").Handler(staticFileHandler).Methods("GET")
	return r
}

func main() {
	r := muxRouter()

	http.ListenAndServe(":8080", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
