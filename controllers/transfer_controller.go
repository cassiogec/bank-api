package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cassiogec/bank-api/api/auth"
	"github.com/cassiogec/bank-api/api/models"
	"github.com/cassiogec/bank-api/api/responses"
)

func (server *Server) AllTransfers(w http.ResponseWriter, r *http.Request) {

	transfer := models.Transfer{}

	transfers, err := transfer.FindAllTransfers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, transfers)
}

func (server *Server) NewTransfer(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	transfer := models.Transfer{}
	err = json.Unmarshal(body, &transfer)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	id, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	transfer.Prepare(id)
	err = transfer.Validate(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	transferCreated, err := transfer.SaveTransfer(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, transferCreated.ID))
	responses.JSON(w, http.StatusCreated, transferCreated)
}
