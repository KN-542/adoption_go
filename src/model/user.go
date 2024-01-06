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
	Name string `json:"name" gorm:"unique;type:varchar(30);index"`
	// メールアドレス
	Email string `json:"email" gorm:"unique;type:varchar(50);check:email ~ '^[a-zA-Z0-9_+-]+(\\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$';index"`
	// パスワード(ハッシュ化)
	Password string `json:"password"`
	// 初回パスワード(ハッシュ化)
	InitPassword string `json:"init_password"`
	// ロールID
	RoleID uint `json:"role_id" gorm:"index"`
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

// ユーザーグループ
type UserGroup struct {
	// ID
	ID uint64 `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	// ハッシュキー
	HashKey string `json:"hash_key" gorm:"unique;type:text"`
	// グループ名
	Name string `json:"name" gorm:"unique;type:varchar(30);index"`
	// 所属ユーザー
	Users string `json:"users" gorm:"type:text;index"`
	// 登録日時
	CreatedAt time.Time `json:"created_at"`
	// 更新日時
	UpdatedAt time.Time `json:"updated_at"`
}

func (t UserGroup) TableName() string {
	return "t_user_group"
}

// ユーザー予定
type UserSchedule struct {
	// ID
	ID uint64 `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	// ハッシュキー
	HashKey string `json:"hash_key" gorm:"unique;type:text"`
	// ハッシュキー(ユーザー)
	UserHashKeys string `json:"user_hash_keys" gorm:"notNull;type:text;index"`
	// タイトル
	Title string `json:"title" gorm:"type:varchar(30)"`
	// 頻度ID
	FreqID uint `json:"freq_id"`
	// 開始時刻
	Start time.Time `json:"start" gorm:"notNull"`
	// 終了時刻
	End time.Time `json:"end" gorm:"notNull"`
	// 登録日時
	CreatedAt time.Time `json:"created_at"`
	// 更新日時
	UpdatedAt time.Time `json:"updated_at"`
	// 頻度(外部キー)
	CalendarFreqStatus CalendarFreqStatus `gorm:"foreignKey:freq_id;references:id"`
}

func (t UserSchedule) TableName() string {
	return "t_user_schedule"
}

// ユーザー予定 Request
type UserScheduleRequest struct {
	// ハッシュキーリスト(ユーザー)
	UserHashKeys string `json:"user_hash_keys"`
	// タイトル
	Title string `json:"title"`
	// 頻度ID
	FreqID uint `json:"freq_id"`
	// 開始時刻
	Start time.Time `json:"start"`
	// 終了時刻
	End time.Time `json:"end"`
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
	// ロール名
	RoleNameJa string `json:"role_name_ja"`
}
type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

// ユーザーグループ response
type UserGroupResponse struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// グループ名
	Name string `json:"name"`
	// 所属ユーザー
	Users string `json:"users"`
}
type UserGroupsResponse struct {
	UserGroups []UserGroupResponse `json:"user_groups"`
}

type CalendarsFreqStatus struct {
	List []CalendarFreqStatus `json:"list"`
}

// ユーザー予定 Response
type UserScheduleResponse struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// ハッシュキー(ユーザー)
	UserHashKeys string `json:"user_hash_keys"`
	// タイトル
	Title string `json:"title"`
	// 開始時刻
	Start time.Time `json:"start"`
	// 終了時刻
	End time.Time `json:"end"`
	// 頻度ID
	FreqID uint `json:"freq_id"`
	// 頻度
	Freq string `json:"freq"`
}

type UserSchedulesResponse struct {
	List []UserScheduleResponse `json:"list"`
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
