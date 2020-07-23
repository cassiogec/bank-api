package main

import (
	"fmt"
	"net/http"

	"bank-api/account"

	"github.com/gorilla/mux"
)

func SetupRoutes(mux *mux.Router) {
	setupGenericRoutes(mux)
	setupAccountRoutes(mux)
	setupTransferRoutes(mux)
	setupLoginRoutes(mux)
}

func setupGenericRoutes(mux *mux.Router) {
	mux.HandleFunc("/", homepage).Methods("GET")
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bank-Api")
}

func setupAccountRoutes(mux *mux.Router) {
	mux.HandleFunc("/accounts", account.AllAccounts).Methods("GET")
	mux.HandleFunc("/user/{account_id}/balance", account.AccountBalance).Methods("GET")
	mux.HandleFunc("/accounts", account.NewAccount).Methods("POST")
}

func setupTransferRoutes(mux *mux.Router) {
}

func setupLoginRoutes(mux *mux.Router) {
}
