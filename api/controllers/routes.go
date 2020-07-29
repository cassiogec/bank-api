package controllers

import (
	"github.com/cassiogec/bank-api/api/middlewares"
	"github.com/urfave/negroni"
)

func (s *Server) initializeRoutes() {

	n := negroni.New(negroni.HandlerFunc(middlewares.JSON), negroni.NewRecovery(), negroni.NewLogger())

	// Home Route
	s.Router.HandleFunc("/", s.Home).Methods("GET")

	//Accounts routes
	s.Router.HandleFunc("/accounts", s.AllAccounts).Methods("GET")
	s.Router.HandleFunc("/accounts/{account_id}/balance", s.AccountBalance).Methods("GET")
	s.Router.HandleFunc("/accounts", s.NewAccount).Methods("POST")

	// Login Route
	s.Router.HandleFunc("/login", s.Login).Methods("POST")

	//Transfers routes
	t := s.Router.PathPrefix("/transfers").Subrouter()
	t.HandleFunc("", s.AllTransfers).Methods("GET")
	t.HandleFunc("", s.NewTransfer).Methods("POST")
	s.Router.PathPrefix("/transfers").Handler(n.With(negroni.HandlerFunc(middlewares.Authentication), negroni.Wrap(t)))

	n.UseHandler(s.Router)
	s.Handler = n
}
