package validator

import (
	"api/src/model/request"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ILoginValidator interface {
	// ログイン
	Login(u *request.Login) error
	// 認証コード生成
	CodeGenerate(u *request.CodeGenerate) error
	// MFA
	MFA(u *request.MFA) error
	// JWT検証
	JWTDecode(u *request.JWTDecode) error
	// パスワード変更
	PasswordChange(u *request.PasswordChange) error
	// ログアウト
	Logout(u *request.Logout) error
}

type LoginValidator struct{}

func NewLoginValidator() ILoginValidator {
	return &LoginValidator{}
}

// ログイン
func (v *LoginValidator) Login(u *request.Login) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Email,
			validation.Required,
			validation.Length(1, 50),
			is.Email,
		),
		validation.Field(
			&u.Password,
			validation.Required,
			validation.Length(8, 16),
			is.Alphanumeric,
		),
	)
}

// 認証コード生成
func (v *LoginValidator) CodeGenerate(u *request.CodeGenerate) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

// MFA
func (v *LoginValidator) MFA(u *request.MFA) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
		validation.Field(
			&u.Code,
			validation.Required,
			validation.Length(6, 6),
			is.UTFNumeric,
		),
	)
}

// JWT検証
func (v *LoginValidator) JWTDecode(u *request.JWTDecode) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

// パスワード変更
func (v *LoginValidator) PasswordChange(u *request.PasswordChange) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
		validation.Field(
			&u.Password,
			validation.Required,
			validation.Length(8, 16),
			is.Alphanumeric,
		),
		validation.Field(
			&u.InitPassword,
			validation.Required,
			validation.Length(8, 16),
			is.Alphanumeric,
			validation.By(func(value interface{}) error {
				initPassword, _ := value.(string)
				if initPassword == u.Password {
					return errors.New("password cannot be the same as the initial password")
				}
				return nil
			}),
		),
	)
}

// ログアウト
func (v *LoginValidator) Logout(u *request.Logout) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}
