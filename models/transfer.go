package models

import (
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

func (t *Transfer) FindAllTransfers(db *gorm.DB) (*[]Transfer, error) {
	var err error
	transfers := []Transfer{}
	err = db.Debug().Model(&Transfer{}).Preload("Account_origin").Preload("Account_destination").Limit(100).Find(&transfers).Error
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
