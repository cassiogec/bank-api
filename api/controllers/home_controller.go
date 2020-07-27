package controllers

import (
	"net/http"

	"github.com/cassiogec/bank-api/api/responses"
)

type Message struct {
	message string
}

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, Message{message: "Welcome to the Bank API"})
}
