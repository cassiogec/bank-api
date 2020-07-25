package main

import (
	"fmt"
	"net/http"

	"github.com/cassiogec/bank-api/controllers"

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
	router.HandleFunc("/accounts", controllers.AllAccounts).Methods("GET")
	router.HandleFunc("/accounts/{account_id}/balance", controllers.AccountBalance).Methods("GET")
	router.HandleFunc("/accounts", controllers.NewAccount).Methods("POST")
}

func setupTransferRoutes(router *mux.Router) {
	router.HandleFunc("/transfers", controllers.AllTransfers).Methods("GET")
	router.HandleFunc("/transfers", controllers.NewTransfer).Methods("POST")
}

func setupLoginRoutes(router *mux.Router) {
	router.HandleFunc("/login", controllers.Login).Methods("POST")
}
