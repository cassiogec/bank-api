package controllers

import (
	"net/http"

	"github.com/cassiogec/bank-api/models"
	"github.com/cassiogec/bank-api/responses"
)

func (server *Server) AllAccounts(w http.ResponseWriter, r *http.Request) {
	account := models.Account{}
	accounts, err := account.FindAllAccounts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, accounts)
}

func AccountBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func NewAccount(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
