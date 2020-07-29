package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/cassiogec/bank-api/api/auth"
	"github.com/cassiogec/bank-api/api/models"
	"github.com/cassiogec/bank-api/api/responses"
	"github.com/jinzhu/gorm"
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
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(account.CPF, account.Secret)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(cpf string, secret string) (auth.Token, error) {

	var err error

	account := models.Account{}

	err = server.DB.Debug().Model(models.Account{}).Where("cpf = ?", cpf).Take(&account).Error
	if gorm.IsRecordNotFoundError(err) {
		return auth.Token{}, errors.New("Account Not Found")
	}
	if err != nil {
		return auth.Token{}, err
	}
	err = models.VerifySecret(account.Secret, secret)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return auth.Token{}, errors.New("Incorrect Secret")
	}
	return auth.CreateToken(account.ID)
}
