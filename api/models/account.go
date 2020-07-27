package models

import (
	"errors"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
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

func (a *Account) FindAccountBalanceByID(db *gorm.DB, id uint64) (float64, error) {
	a, err := a.FindAccountByID(db, id)
	if err != nil {
		return 0, err
	}
	return a.Balance, nil
}

func (a *Account) FindAccountByID(db *gorm.DB, id uint64) (*Account, error) {
	var err error
	err = db.Debug().Model(Account{}).Where("id = ?", id).Take(&a).Error
	if gorm.IsRecordNotFoundError(err) {
		return &Account{}, errors.New("Account Not Found")
	}
	if err != nil {
		return &Account{}, err
	}
	a.SanitizeAccount()
	return a, nil
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

func (a *Account) Prepare() {
	a.ID = 0
	a.Name = html.EscapeString(strings.TrimSpace(a.Name))
	a.CPF = html.EscapeString(strings.TrimSpace(a.CPF))
	a.Balance = 10
	a.CreatedAt = time.Now()
}

func (a *Account) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if err := a.ValidateSecret(); err != nil {
			return err
		}
		if err := a.ValidateCPF(); err != nil {
			return err
		}
		return nil

	default:
		if err := a.ValidateName(); err != nil {
			return err
		}
		if err := a.ValidateSecret(); err != nil {
			return err
		}
		if err := a.ValidateCPF(); err != nil {
			return err
		}
		return nil
	}
}

func (a *Account) ValidateName() error {
	if a.Name == "" {
		return errors.New("Required Name")
	}
	return nil
}

func (a *Account) ValidateSecret() error {
	if a.Secret == "" {
		return errors.New("Required Secret")
	}
	return nil
}

func (a *Account) ValidateCPF() error {
	if a.CPF == "" {
		return errors.New("Required CPF")
	}
	if len(a.CPF) != 11 {
		return errors.New("Invalid CPF")
	}
	if _, err := strconv.Atoi(a.CPF); err != nil {
		return errors.New("Invalid CPF")
	}
	return nil
}

func (a *Account) SaveAccount(db *gorm.DB) (*Account, error) {
	var err error
	if err := a.ValidateUniqueCPF(db); err != nil {
		return &Account{}, err
	}
	err = db.Debug().Create(&a).Error
	if err != nil {
		return &Account{}, err
	}
	a.SanitizeAccount()
	return a, nil
}

func (a *Account) ValidateUniqueCPF(db *gorm.DB) error {
	accountFound, _ := a.FindAccountByCPF(db, a.CPF)
	if accountFound.ID != 0 {
		return errors.New("There is already an account with the given CPF")
	}
	return nil
}

func (a *Account) FindAccountByCPF(db *gorm.DB, cpf string) (*Account, error) {
	var err error
	err = db.Debug().Model(Account{}).Where("cpf = ?", cpf).Take(&a).Error
	if err != nil {
		return &Account{}, err
	}
	return a, nil
}

func (a *Account) BeforeSave() error {
	hashedSecret, err := Hash(a.Secret)
	if err != nil {
		return err
	}
	a.Secret = string(hashedSecret)
	return nil
}

func Hash(secret string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
}

func VerifySecret(hashedSecret, secret string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(secret))
}

func (a *Account) WithdrawFromBalance(db *gorm.DB, amount float64) error {
	newBalance := a.Balance - amount
	return a.UpdateBalance(db, newBalance)
}

func (a *Account) DepositOnBalance(db *gorm.DB, amount float64) error {
	newBalance := a.Balance + amount
	return a.UpdateBalance(db, newBalance)
}

func (a *Account) UpdateBalance(db *gorm.DB, balance float64) error {
	db = db.Debug().Model(&Account{}).Where("id = ?", a.ID).Take(&a).UpdateColumns(
		map[string]interface{}{
			"balance": balance,
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
