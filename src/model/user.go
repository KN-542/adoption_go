package model

import (
	"time"
)

// ユーザー(管理)
type User struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
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

// ユーザー(管理) response
type UserResponse struct {
	// ID
	ID uint `json:"id"`
	// 氏名
	Name string `json:"name"`
	// メールアドレス
	Email string `json:"email"`
	// ロールID
	RoleID uint `json:"role_id"`
}

func ConvertUser(u *[]User) *[]UserResponse {
	var respList []UserResponse
	for _, row := range *u {
		respList = append(
			respList,
			UserResponse{
				ID:     row.ID,
				Name:   row.Name,
				Email:  row.Email,
				RoleID: row.RoleID,
			},
		)
	}
	return &respList
}
