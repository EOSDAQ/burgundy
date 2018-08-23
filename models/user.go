package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type secretString string

// User ...
type User struct {
	gorm.Model   `json:"-"`
	AccountName  string        `json:"accountName" gorm:"not null;unique"`
	Email        string        `json:"email"`
	EmailHash    *secretString `json:"emailHash,omitempty"`
	EmailConfirm bool          `json:"emailConfirm"`
	OTPKey       *secretString `json:"otpKey,omitempty""`
	OTPConfirm   bool          `json:"otpConfirm"`
	Registered   bool          `json:"-"`
}

func (u *User) String() string {
	return fmt.Sprintf("AccountName[%s] Email[%s]", u.AccountName, u.Email)
}

func (u *User) UpdateConfirm(other *User) {
	if other.EmailConfirm != u.EmailConfirm {
		if other.EmailConfirm && other.EmailHash == u.EmailHash {
			u.EmailConfirm = other.EmailConfirm
		} else {
			u.EmailConfirm = other.EmailConfirm
		}
	}
	if other.OTPConfirm != u.OTPConfirm {
		if other.OTPConfirm && other.OTPKey == u.OTPKey {
			u.OTPConfirm = other.OTPConfirm
		} else {
			u.OTPConfirm = other.OTPConfirm
		}
		u.OTPKey = other.OTPKey
	}
}

func (u *User) NeedRegister() bool {
	mlog.Debugw("NeedRegister", "email", u.EmailConfirm, "otp", u.OTPConfirm, "r", u.Registered)
	return u.EmailConfirm && u.OTPConfirm && !u.Registered
}

func (u *User) NeedUnregister() bool {
	mlog.Debugw("NeedUnrgister", "user", u)
	return !(u.EmailConfirm && u.OTPConfirm) && u.Registered
}

func (u *User) UpdateRegister() {
	u.Registered = !u.Registered
}
