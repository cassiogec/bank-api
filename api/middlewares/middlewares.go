package middlewares

import (
	"errors"
	"net/http"

	"github.com/cassiogec/bank-api/api/auth"
	"github.com/cassiogec/bank-api/api/responses"
)

func JSON(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Set("Content-Type", "application/json")
	next(rw, r)
}

func Authentication(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := auth.TokenValid(r)
	if err != nil {
		responses.ERROR(rw, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	next(rw, r)
}
