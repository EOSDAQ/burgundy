package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// User ...
type User struct {
	gorm.Model   `json:"-"`
	AccountName  string `json:"account_name"`
	Email        string `json:"email"`
	ConfirmEmail bool   `json:"confirm_email"`
	ConfirmOTP   bool   `json:"confirm_otp"`
	EOSPubkey    string `json:"eos_pubkey"`
}

func (u *User) String() string {
	return fmt.Sprintf("AccountName[%s] Email[%s]", u.AccountName, u.Email)
}
