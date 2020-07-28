package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSignIn(t *testing.T) {

	err := refreshAccountTable()
	if err != nil {
		log.Fatal(err)
	}

	account, err := seedOneAccount()
	if err != nil {
		fmt.Printf("This is the error %v\n", err)
	}

	samples := []struct {
		cpf          string
		secret       string
		errorMessage string
	}{
		{
			cpf:          account.CPF,
			secret:       "doctor10",
			errorMessage: "",
		},
		{
			cpf:          account.CPF,
			secret:       "Wrong secret",
			errorMessage: "Incorrect Secret",
		},
		{
			cpf:          "Wrong CPF",
			secret:       "secret",
			errorMessage: "record not found",
		},
	}

	for _, v := range samples {

		token, err := server.SignIn(v.cpf, v.secret)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

func TestLogin(t *testing.T) {

	refreshAccountTable()

	_, err := seedOneAccount()
	if err != nil {
		fmt.Printf("This is the error %v\n", err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		cpf          string
		secret       string
		errorMessage string
	}{
		{
			inputJSON:    `{"cpf": "11111111111", "secret": "doctor10"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    `{"cpf": "11111111111", "secret": "wrong secret"}`,
			statusCode:   422,
			errorMessage: "Incorrect Secret",
		},
		{
			inputJSON:    `{"cpf": "sdf4565", "secret": "secret"}`,
			statusCode:   422,
			errorMessage: "Invalid CPF",
		},
		{
			inputJSON:    `{"cpf": "", "secret": "secret"}`,
			statusCode:   422,
			errorMessage: "Required CPF",
		},
		{
			inputJSON:    `{"cpf": "11111111111", "secret": ""}`,
			statusCode:   422,
			errorMessage: "Required Secret",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}

		if v.statusCode == 422 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
