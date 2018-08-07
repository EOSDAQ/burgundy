package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// User ...
type User struct {
	gorm.Model   `json:"-"`
	AccountName  string `json:"accountName"`
	Email        string `json:"email"`
	EmailHash    string `json:"emailHash"`
	EmailConfirm bool   `json:"emailConfirm"`
	OTPKey       string `json:"otpKey"`
	OTPConfirm   bool   `json:"otpConfirm"`
}

func (u *User) String() string {
	return fmt.Sprintf("AccountName[%s] Email[%s]", u.AccountName, u.Email)
}
