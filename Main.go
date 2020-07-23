package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handlers struct{}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	handleRoutes(myRouter)
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func handleRoutes(mux *mux.Router) {

}
