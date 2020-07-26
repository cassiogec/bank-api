package controllers

import (
	"net/http"

	"github.com/cassiogec/bank-api/models"
	"github.com/cassiogec/bank-api/responses"
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

func NewTransfer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
