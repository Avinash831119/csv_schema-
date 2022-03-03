package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func initializeRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/schema", GetSchema).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	initializeRouter()
}
