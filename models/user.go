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
	EmailHash    string `json:"email_hash"`
	EmailConfirm bool   `json:"email_confirm"`
	OTPKey       string `json:"otp_key"`
	OTPConfirm   bool   `json:"otp_confirm"`
}

func (u *User) String() string {
	return fmt.Sprintf("AccountName[%s] Email[%s]", u.AccountName, u.Email)
}
