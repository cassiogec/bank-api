package tests

import (
	"log"
	"testing"

	"github.com/cassiogec/bank-api/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllTransfers(t *testing.T) {

	err := refreshAccountAndTransferTable()
	if err != nil {
		log.Fatalf("Error refreshing account and transfer table %v\n", err)
	}
	_, _, err = seedAccountsAndTransfers()
	if err != nil {
		log.Fatalf("Error seeding account and transfer table %v\n", err)
	}
	transfers, err := transferInstance.FindAllTransfers(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the transfers: %v\n", err)
		return
	}
	assert.Equal(t, len(*transfers), 2)
}

func TestSaveTransfer(t *testing.T) {

	err := refreshAccountAndTransferTable()
	if err != nil {
		log.Fatalf("Error account and transfer refreshing table %v\n", err)
	}

	accounts, err := seedAccounts()
	if err != nil {
		log.Fatalf("Cannot seed accounts %v\n", err)
	}

	newTransfer := models.Transfer{
		ID:                     1,
		Account_origin_id:      accounts[0].ID,
		Account_origin:         accounts[0],
		Account_destination_id: accounts[1].ID,
		Account_destination:    accounts[1],
		Amount:                 5,
	}
	savedTransfer, err := newTransfer.SaveTransfer(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the transfer: %v\n", err)
		return
	}
	assert.Equal(t, newTransfer.ID, savedTransfer.ID)
	assert.Equal(t, newTransfer.Account_origin_id, savedTransfer.Account_origin_id)
	assert.Equal(t, newTransfer.Account_destination_id, savedTransfer.Account_destination_id)
	assert.Equal(t, newTransfer.Amount, savedTransfer.Amount)
}
