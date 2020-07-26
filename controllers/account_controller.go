package controllers

import (
	"net/http"
	"strconv"

	"github.com/cassiogec/bank-api/models"
	"github.com/cassiogec/bank-api/responses"
	"github.com/gorilla/mux"
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

func (server *Server) AccountBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["account_id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	account := models.Account{}
	balance, err := account.FindAccountBalanceByID(server.DB, uint64(id))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, balance)
}

func NewAccount(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
