package controllers

import (
	"net/http"

	"github.com/cassiogec/bank-api/api/responses"
)

type Message struct {
	Message string `json:"message"`
}

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	message := Message{Message: "Welcome to the Bank API"}
	responses.JSON(w, http.StatusOK, message)
}
