package tests

import (
	"log"
	"testing"

	"github.com/cassiogec/bank-api/api/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllAccounts(t *testing.T) {

	err := refreshAccountTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedAccounts()
	if err != nil {
		log.Fatal(err)
	}

	accounts, err := accountInstance.FindAllAccounts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the accounts: %v\n", err)
		return
	}
	assert.Equal(t, len(*accounts), 2)
}

func TestSaveAccount(t *testing.T) {

	err := refreshAccountTable()
	if err != nil {
		log.Fatal(err)
	}

	newAccount := models.Account{
		ID:     1,
		CPF:    "00000000000",
		Name:   "test",
		Secret: "secret",
	}
	savedAccount, err := newAccount.SaveAccount(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the accounts: %v\n", err)
		return
	}
	assert.Equal(t, newAccount.ID, savedAccount.ID)
	assert.Equal(t, newAccount.CPF, savedAccount.CPF)
	assert.Equal(t, newAccount.Name, savedAccount.Name)
}

func TestFindAccountByID(t *testing.T) {

	err := refreshAccountTable()
	if err != nil {
		log.Fatal(err)
	}

	account, err := seedOneAccount()
	if err != nil {
		log.Fatalf("cannot seed accounts table: %v", err)
	}
	foundAccount, err := accountInstance.FindAccountByID(server.DB, account.ID)
	if err != nil {
		t.Errorf("this is the error getting one account: %v\n", err)
		return
	}
	assert.Equal(t, foundAccount.ID, account.ID)
	assert.Equal(t, foundAccount.CPF, account.CPF)
	assert.Equal(t, foundAccount.Name, account.Name)
}
