package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type secretString = string

// User ...
type User struct {
	gorm.Model  `json:"-"`
	AccountName string        `json:"accountName" gorm:"not null;unique"`
	AccountHash *secretString `json:"accountHash,omitempty" gorm:"not null"`

	Email        string        `json:"email"`
	EmailHash    *secretString `json:"emailHash,omitempty"`
	EmailConfirm bool          `json:"emailConfirm"`

	OTPKey     string `json:"-"`
	OTPConfirm bool   `json:"otpConfirm"`
	Registered bool   `json:"-"`
}

func (u *User) String() string {
	return fmt.Sprintf("AccountName[%s] Email[%s]", u.AccountName, u.Email)
}

func (u *User) Login(hash string) bool {
	return string(*u.AccountHash) == hash
}

func (u *User) ConfirmEmail(email, emailHash string) bool {
	if u.Email == email && string(*u.EmailHash) == emailHash {
		u.EmailConfirm = true
	}
	return u.EmailConfirm
}

func (u *User) RevokeEmail(email, emailHash string) {
	if email != "" {
		u.Email = email
	}
	u.EmailHash = &emailHash
	u.EmailConfirm = false
}

func (u *User) Validate() bool {
	return u.AccountName != "" && u.ID == 0 &&
		u.Email != "" && u.EmailHash != nil && *u.EmailHash != "" &&
		u.EmailConfirm == false && u.OTPConfirm == false
}

func (u *User) NeedRegister() bool {
	return u.EmailConfirm && u.OTPConfirm && !u.Registered
}

func (u *User) NeedUnregister() bool {
	return !(u.EmailConfirm && u.OTPConfirm) && u.Registered
}

func (u *User) UpdateRegister() {
	u.Registered = !u.Registered
}

func (u *User) GenerateOTPKey() (string, error) {
	if u.OTPKey != "" {
		return "", fmt.Errorf("Already exists OTP Key [%s]", u.AccountName)
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "eosdaq.com",
		AccountName: u.AccountName,
	})
	if err != nil {
		return "", fmt.Errorf("GenerateOTPKey account[%s] error[%s]", u.AccountName, err)
	}
	u.OTPKey = key.Secret()
	return u.OTPKey, nil
}

func (u *User) RemoveOTPKey() {
	u.OTPKey = ""
	u.OTPConfirm = false
	return
}

func (u *User) ValidateOTP(code string) (ok bool) {
	if u.OTPKey == "" {
		return false
	}

	keyURL := fmt.Sprintf("otpauth://totp/eosdaq.com:%s?secret=%s&issuer=eosdaq.com&algorithm=SHA1&digits=6&period=30", u.AccountName, u.OTPKey)
	key, err := otp.NewKeyFromURL(keyURL)
	if err != nil {
		mlog.Infow("ValidateOTP error", "account", u.AccountName, "err", err)
		return false
	}

	ok = totp.Validate(code, key.Secret())
	if ok && !u.OTPConfirm {
		u.OTPConfirm = true
	}
	return ok
}
