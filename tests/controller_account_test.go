package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/cassiogec/bank-api/api/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestNewAccount(t *testing.T) {

	err := refreshAccountTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		name         string
		cpf          string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"David Tennant", "cpf": "11111111111", "secret": "doctor10"}`,
			statusCode:   201,
			name:         "David Tennant",
			cpf:          "11111111111",
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"David", "cpf": "11111111111", "secret": "secret"}`,
			statusCode:   500,
			errorMessage: "CPF Already Taken",
		},
		{
			inputJSON:    `{"name":"David", "cpf": "22ff8981a", "secret": "secret"}`,
			statusCode:   422,
			errorMessage: "Invalid CPF",
		},
		{
			inputJSON:    `{"name": "", "cpf": "11111111111", "secret": "secret"}`,
			statusCode:   422,
			errorMessage: "Required Name",
		},
		{
			inputJSON:    `{"name": "David", "cpf": "", "secret": "secret"}`,
			statusCode:   422,
			errorMessage: "Required CPF",
		},
		{
			inputJSON:    `{"name": "Kan", "cpf": "kan@gmail.com", "secret": ""}`,
			statusCode:   422,
			errorMessage: "Required Secret",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/accounts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.NewAccount)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["cpf"], v.cpf)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestAllAccounts(t *testing.T) {

	err := refreshAccountTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedAccounts()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/accounts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.AllAccounts)
	handler.ServeHTTP(rr, req)

	var accounts []models.Account
	err = json.Unmarshal([]byte(rr.Body.String()), &accounts)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(accounts), 2)
}

func TestFindAccountBalanceByID(t *testing.T) {

	err := refreshAccountTable()
	if err != nil {
		log.Fatal(err)
	}
	account, err := seedOneAccount()
	if err != nil {
		log.Fatal(err)
	}
	accountSample := []struct {
		id           string
		statusCode   int
		name         string
		cpf          string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(account.ID)),
			statusCode: 200,
			name:       account.Name,
			cpf:        account.CPF,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range accountSample {

		req, err := http.NewRequest("GET", "/accounts", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"account_id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.AccountBalance)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, account.Balance, responseMap["Balance"])
		}
	}
}
