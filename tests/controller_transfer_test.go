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
	"gopkg.in/go-playground/assert.v1"
)

func TestNewTransfer(t *testing.T) {

	err := refreshAccountAndTransferTable()
	if err != nil {
		log.Fatal(err)
	}
	accounts, err := seedAccounts()
	if err != nil {
		log.Fatalf("Cannot seed accounts %v\n", err)
	}
	token, err := server.SignIn(accounts[0].CPF, accounts[0].Secret)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}

	samples := []struct {
		inputJSON              string
		statusCode             int
		account_origin_id      uint64
		account_destination_id uint64
		amount                 float64
		tokenGiven             string
		errorMessage           string
	}{
		{
			inputJSON:              `{"account_destination_id": ` + strconv.FormatUint(accounts[1].ID, 10) + `, "amount": 3.5}`,
			statusCode:             201,
			tokenGiven:             token.Token,
			account_origin_id:      accounts[0].ID,
			account_destination_id: accounts[1].ID,
			amount:                 3.5,
			errorMessage:           "",
		},
		{
			inputJSON:    `{"account_destination_id": ` + strconv.FormatUint(accounts[1].ID, 10) + `, "amount": 3.5}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"account_destination_id": ` + strconv.FormatUint(accounts[1].ID, 10) + `, "amount": 3.5}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"amount": 3.5}`,
			statusCode:   422,
			tokenGiven:   token.Token,
			errorMessage: "Required Account Destination ID",
		},
		{
			inputJSON:    `{"account_destination_id": ` + strconv.FormatUint(accounts[0].ID, 10) + `, "amount": 3.5}`,
			statusCode:   422,
			tokenGiven:   token.Token,
			errorMessage: "Origin and Destination Accouts should differ",
		},
		{
			inputJSON:    `{"account_destination_id": ` + strconv.FormatUint(accounts[1].ID, 10) + `, "amount": 1000000000}`,
			statusCode:   422,
			tokenGiven:   token.Token,
			errorMessage: "Insufficient balance",
		},
		{
			inputJSON:    `{"account_destination_id": ` + strconv.FormatUint(accounts[1].ID, 10) + `, "amount": -3.5}`,
			statusCode:   422,
			tokenGiven:   token.Token,
			errorMessage: "Amount should be bigger then 0",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/transfers", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.NewTransfer)

		req.Header.Set("Authorization", "token: "+v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["account_origin_id"], float64(v.account_origin_id))           //the response id is float64
			assert.Equal(t, responseMap["account_destination_id"], float64(v.account_destination_id)) //the response id is float64
			assert.Equal(t, responseMap["amount"], v.amount)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestAllTransfers(t *testing.T) {

	err := refreshAccountAndTransferTable()
	if err != nil {
		log.Fatal(err)
	}
	accounts, _, err := seedAccountsAndTransfers()
	if err != nil {
		log.Fatal(err)
	}

	token, err := server.SignIn(accounts[0].CPF, accounts[0].Secret)
	if err != nil {
		log.Fatalf("cannot login: %v\n", accounts[0].Secret)
	}

	req, err := http.NewRequest("GET", "/transfers", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	req.Header.Set("Authorization", "token: "+token.Token)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.AllTransfers)
	handler.ServeHTTP(rr, req)

	var transfers []models.Transfer
	err = json.Unmarshal([]byte(rr.Body.String()), &transfers)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(transfers), 2)
}
