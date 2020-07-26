package controllers

import "github.com/cassiogec/bank-api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Accounts routes
	s.Router.HandleFunc("/accounts", middlewares.SetMiddlewareJSON(s.AllAccounts)).Methods("GET")
	s.Router.HandleFunc("/accounts/{account_id}/balance", middlewares.SetMiddlewareJSON(s.AccountBalance)).Methods("GET")
	s.Router.HandleFunc("/accounts", middlewares.SetMiddlewareJSON(s.NewAccount)).Methods("POST")
}
