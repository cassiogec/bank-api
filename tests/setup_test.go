package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/cassiogec/bank-api/api/controllers"
	"github.com/cassiogec/bank-api/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var accountInstance = models.Account{}
var transferInstance = models.Transfer{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TEST_DB_DRIVER")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PASSWORD"))
	server.DB, err = gorm.Open(TestDbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to the database\n")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the database\n")
	}
}

func refreshAccountTable() error {
	err := server.DB.DropTableIfExists(&models.Account{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Account{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneAccount() (models.Account, error) {

	err := refreshAccountTable()
	if err != nil {
		log.Fatal(err)
	}

	account := models.Account{
		Name:   "David Tennant",
		CPF:    "11111111111",
		Secret: "doctor10",
	}

	err = server.DB.Model(&models.Account{}).Create(&account).Error
	if err != nil {
		return models.Account{}, err
	}
	return account, nil
}

func seedAccounts() ([]models.Account, error) {
	accounts := []models.Account{
		models.Account{
			Name:   "Tom Baker",
			CPF:    "22222222222",
			Secret: "doctor4",
		},
		models.Account{
			Name:   "Matt Smith",
			CPF:    "33333333333",
			Secret: "doctor11",
		},
	}

	newAccounts := make([]models.Account, 2)
	for i, account := range accounts {
		err := server.DB.Model(&models.Account{}).Create(&account).Error
		newAccounts[i] = account
		if err != nil {
			return []models.Account{}, err
		}
	}
	return newAccounts, nil
}

func refreshAccountAndTransferTable() error {

	err := server.DB.DropTableIfExists(&models.Account{}, &models.Transfer{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Account{}, &models.Transfer{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedAccountsAndTransfers() ([]models.Account, []models.Transfer, error) {

	var err error

	if err != nil {
		return []models.Account{}, []models.Transfer{}, err
	}
	var accounts = []models.Account{
		models.Account{
			Name:   "Peter Davison",
			CPF:    "55555555555",
			Secret: "doctor5",
		},
		models.Account{
			Name:   "John Hurt",
			CPF:    "66666666666",
			Secret: "wardoctor",
		},
	}

	for i, _ := range accounts {
		err = server.DB.Model(&models.Account{}).Create(&accounts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed accounts table: %v", err)
		}
	}

	var transfers = []models.Transfer{
		models.Transfer{
			Account_origin_id:      accounts[0].ID,
			Account_destination_id: accounts[1].ID,
			Amount:                 2.5,
		},
		models.Transfer{
			Account_origin_id:      accounts[1].ID,
			Account_destination_id: accounts[0].ID,
			Amount:                 7,
		},
	}

	for i, _ := range transfers {
		err = server.DB.Model(&models.Transfer{}).Create(&transfers[i]).Error
		if err != nil {
			log.Fatalf("cannot seed transfers table: %v", err)
		}
	}

	return accounts, transfers, nil
}
