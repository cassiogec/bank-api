package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Account struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;" json:"name"`
	CPF       string    `gorm:"size:11; not null;unique" json:"cpf"`
	Secret    string    `gorm:"size:100;not null;" json:"secret"`
	Balance   float64   `gorm:"default:0; not null" json:"balance"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (a *Account) FindAllAccounts(db *gorm.DB) (*[]Account, error) {
	var err error
	accounts := []Account{}
	err = db.Debug().Model(&Account{}).Limit(100).Find(&accounts).Error
	if err != nil {
		return &[]Account{}, err
	}
	SanitizeAccounts(&accounts)
	return &accounts, err
}

func SanitizeAccounts(accounts *[]Account) {
	for i, a := range *accounts {
		a.SanitizeAccount()
		(*accounts)[i] = a
	}
}

func (a *Account) SanitizeAccount() {
	a.Secret = ""
}
