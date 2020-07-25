package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	SetupRoutes(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}
