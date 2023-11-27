package model

import (
	"time"
)

// ユーザ(管理)
type User struct {
	// ID
	ID uint64 `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	// ハッシュキー
	HashKey string `json:"hash_key" gorm:"unique;type:text"`
	// 氏名
	Name string `json:"name" gorm:"unique;type:varchar(30)"`
	// メールアドレス
	Email string `json:"email" gorm:"unique;type:varchar(50);check:email ~ '^[a-zA-Z0-9_+-]+(\\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$'"`
	// パスワード(ハッシュ化)
	Password string `json:"password"`
	// 初回パスワード(ハッシュ化)
	InitPassword string `json:"init_password"`
	// ロールID
	RoleID uint `json:"role_id"`
	// 登録日時
	CreatedAt time.Time `json:"created_at"`
	// 更新日時
	UpdatedAt time.Time `json:"updated_at"`
	// ロール(外部キー)
	Role Role `gorm:"foreignKey:role_id;references:id"`
}

func (t User) TableName() string {
	return "t_user"
}

// ユーザ(管理) response
type UserResponse struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// 氏名
	Name string `json:"name"`
	// メールアドレス
	Email string `json:"email"`
	// ロールID
	RoleID uint `json:"role_id"`
	// 初回パスワード
	InitPassword string `json:"init_password"`
	// MFA認証フラグ
	MFA int8 `json:"mfa"`
	// パスワード変更 必要性
	PasswordChange int8 `json:"password_change"`
}
type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

// ユーザー MFA
type UserMFA struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// 認証コード
	Code string `json:"code"`
}

// パスワード変更
type PasswordChanging struct {
	// 旧パスワード(ハッシュ化)
	OldPassword string `json:"old_password"`
	// 新パスワード(ハッシュ化)
	Password string `json:"password"`
}

// ユーザーロール一覧
type UserRoles struct {
	Roles []Role `json:"roles"`
}

func ConvertUsers(u *[]User) *[]UserResponse {
	var respList []UserResponse
	for _, row := range *u {
		respList = append(
			respList,
			UserResponse{
				HashKey: row.HashKey,
				Name:    row.Name,
				Email:   row.Email,
				RoleID:  row.RoleID,
			},
		)
	}
	return &respList
}

func ConvertUser(u *User) *UserResponse {
	return &UserResponse{
		HashKey: u.HashKey,
		Name:    u.Name,
		Email:   u.Email,
		RoleID:  u.RoleID,
	}
}
