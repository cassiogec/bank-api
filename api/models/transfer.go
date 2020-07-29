package models

import (
	"errors"
	"math"
	"time"

	"github.com/jinzhu/gorm"
)

type Transfer struct {
	ID                     uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Account_origin_id      uint64    `gorm:"not null" json:"account_origin_id"`
	Account_origin         Account   `gorm:"foreignkey:Account_origin_id"`
	Account_destination_id uint64    `gorm:"not null" json:"account_destination_id"`
	Account_destination    Account   `gorm:"foreignkey:Account_destination_id"`
	Amount                 float64   `gorm:"default:0; not null" json:"amount"`
	CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (t *Transfer) FindAllTransfersById(db *gorm.DB, account_id uint64) (*[]Transfer, error) {
	var err error
	transfers := []Transfer{}
	err = db.Debug().Model(&Transfer{}).Preload("Account_origin").Preload("Account_destination").Where("account_origin_id = ? OR account_destination_id = ?", account_id, account_id).Find(&transfers).Error
	if err != nil {
		return &[]Transfer{}, err
	}
	SanitizeTransferAccounts(&transfers)
	return &transfers, nil
}

func SanitizeTransferAccounts(transfers *[]Transfer) {
	for i, t := range *transfers {
		t.Account_destination.SanitizeAccount()
		t.Account_origin.SanitizeAccount()
		(*transfers)[i] = t
	}
}

func (t *Transfer) Prepare(account_origin_id uint64) {
	t.ID = 0
	t.Account_origin_id = account_origin_id
	t.Account_origin = Account{}
	t.Account_destination = Account{}
	t.CreatedAt = time.Now()
	t.Amount = math.Round(t.Amount*100) / 100
}

func (t *Transfer) Validate(db *gorm.DB) error {
	if t.Account_destination_id < 1 {
		return errors.New("Required Account Destination ID")
	}
	if t.Account_destination_id == t.Account_origin_id {
		return errors.New("Origin and Destination Accouts should differ")
	}
	if t.Amount <= 0.0 {
		return errors.New("Amount should be bigger then 0")
	}
	if _, err := t.Account_destination.FindAccountByID(db, t.Account_destination_id); err != nil {
		return err
	}
	if _, err := t.Account_origin.FindAccountByID(db, t.Account_origin_id); err != nil {
		return err
	}
	if t.Account_origin.Balance-t.Amount < 0 {
		return errors.New("Insufficient balance")
	}
	return nil
}

func (t *Transfer) SaveTransfer(db *gorm.DB) (*Transfer, error) {
	var err error
	tx := db.Begin()
	err = t.Account_origin.WithdrawFromBalance(db, t.Amount)
	if err != nil {
		tx.Rollback()
		return &Transfer{}, err
	}
	err = t.Account_destination.DepositOnBalance(db, t.Amount)
	if err != nil {
		tx.Rollback()
		return &Transfer{}, err
	}
	err = db.Debug().Model(&Transfer{}).Create(&t).Error
	if err != nil {
		tx.Rollback()
		return &Transfer{}, err
	}
	t.Account_origin.SanitizeAccount()
	t.Account_destination.SanitizeAccount()
	tx.Commit()
	return t, nil
}
