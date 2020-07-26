package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cassiogec/bank-api/auth"
	"github.com/cassiogec/bank-api/models"
	"github.com/cassiogec/bank-api/responses"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	account := models.Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	account.Prepare()
	err = account.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	token, err := server.SignIn(account.CPF, account.Secret)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(cpf string, secret string) (string, error) {

	var err error

	account := models.Account{}

	err = server.DB.Debug().Model(models.Account{}).Where("cpf = ?", cpf).Take(&account).Error
	if err != nil {
		return "", err
	}
	err = models.VerifySecret(account.Secret, secret)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(account.ID)
}
