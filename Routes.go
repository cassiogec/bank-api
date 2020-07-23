package main

import (
	"fmt"
	"net/http"

	"bank-api/account"
	"bank-api/login"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	setupGenericRoutes(router)
	setupAccountRoutes(router)
	setupTransferRoutes(router)
	setupLoginRoutes(router)
}

func setupGenericRoutes(router *mux.Router) {
	router.HandleFunc("/", homepage).Methods("GET")
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bank-Api")
}

func setupAccountRoutes(router *mux.Router) {
	router.HandleFunc("/accounts", account.AllAccounts).Methods("GET")
	router.HandleFunc("/accounts/{account_id}/balance", account.AccountBalance).Methods("GET")
	router.HandleFunc("/accounts", account.NewAccount).Methods("POST")
}

func setupTransferRoutes(mux *mux.Router) {
}

func setupLoginRoutes(router *mux.Router) {
	router.HandleFunc("/login", login.Login).Methods("POST")
}
