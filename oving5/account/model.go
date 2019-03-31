package account

import (
	"fmt"
	"oving5/config"
	"time"
)

// Account ...
type Account struct {
	Nr      string `gorm:"primary_key"`
	Balance float64
	Holder  string
}

// Withdraw ...
func (a *Account) Withdraw(amount float64) {
	var balance struct{ Balance float64 }
	config.DB.
		Table("accounts").
		Where("nr = ?", a.Nr).
		Select("balance").
		First(&balance)

	fmt.Println(balance)
	newBalance := balance.Balance - amount

	config.DB.
		Model(a).
		Where("nr = ?", a.Nr).
		Update("balance", newBalance)
}

// WithdrawSafely ...
func (a *Account) WithdrawSafely(amount float64) {
	var balance struct{ Balance float64 }
	config.DB.
		Table("accounts").
		Where("nr = ?", a.Nr).
		Select("balance").
		First(&balance)

	fmt.Println(balance)
	newBalance := balance.Balance - amount

	affectedRows := config.DB.
		Model(a).
		Where("nr = ? AND balance = ?", a.Nr, balance.Balance).
		Updates(map[string]interface{}{"balance": newBalance}).
		RowsAffected

	if affectedRows == 0 {
		time.Sleep(200 * time.Millisecond)
		a.WithdrawSafely(amount)
	}
}

// Sync ...
func Sync(a *Account) {
	config.DB.Where(&a).First(&a)
}

// FindWhere ...
func FindWhere(where ...interface{}) []Account {
	var users []Account
	//config.DB.Find(users, where)
	config.DB.
		Where(where[0], where[1:]...).
		Order("balance desc").
		Find(&users)
	return users
}

// CreateAccount should be used to create a new account as it will put the new account in the db
func CreateAccount(nr, holder string, balance float64) *Account {
	newAccount := &Account{nr, balance, holder}
	config.DB.
		Where(Account{Nr: nr}).
		Assign(Account{Holder: holder, Balance: balance}).
		FirstOrCreate(&newAccount)
	return newAccount
}

func init() {
	config.DB.AutoMigrate(&Account{})
}
