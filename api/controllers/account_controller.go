package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cassiogec/bank-api/api/models"
	"github.com/cassiogec/bank-api/api/responses"
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
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.JSON(w, http.StatusOK, balance)
}

func (server *Server) NewAccount(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	account := models.Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	account.Prepare()
	err = account.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	accountCreated, err := account.SaveAccount(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, accountCreated.ID))
	responses.JSON(w, http.StatusCreated, accountCreated)
}
