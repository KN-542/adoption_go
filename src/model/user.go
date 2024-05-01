package model

import "time"

/*
	t_user
	ユーザー
*/
type User struct {
	AbstractTransactionModel
	// 氏名
	Name string `json:"name" gorm:"unique;not null;check:name <> '';type:varchar(75);index"`
	// メールアドレス
	Email string `json:"email" gorm:"unique;not null;type:varchar(100);check:email ~ '^[a-zA-Z0-9_+-]+(\\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$';index"`
	// パスワード(ハッシュ化)
	Password string `json:"password" gorm:"not null;check:password <> ''"`
	// 初回パスワード(ハッシュ化)
	InitPassword string `json:"init_password" gorm:"not null;check:init_password <> ''"`
	// ロールID
	RoleID uint `json:"role_id"`
	// ユーザー種別
	UserType uint `json:"user_type"`
	// リフレッシュトークン
	RefreshToken string `json:"refresh_token" gorm:"type:text"`
	// ロール(外部キー)
	Role CustomRole `gorm:"foreignKey:role_id;references:id"`
	// ログイン種別(外部キー)
	LoginType LoginType `gorm:"foreignKey:user_type;references:id"`
}

/*
	t_user_group
	ユーザーグループ
*/
type UserGroup struct {
	AbstractTransactionModel
	// グループ名
	Name string `json:"name" gorm:"not null;unique;check:name <> '';type:varchar(30);index"`
}

/*
	t_user_group_association
	ユーザーグループ紐づけ
*/
type UserGroupAssociation struct {
	// ユーザーグループID
	UserGroupID uint `json:"user_group_id" gorm:"primaryKey"`
	// ユーザーID
	UserID uint `json:"user_id" gorm:"primaryKey"`
	// ユーザーグループ(外部キー)
	UserGroup UserGroup `gorm:"foreignKey:user_group_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

/*
	t_user_schedule
	ユーザー予定
*/
type UserSchedule struct {
	AbstractTransactionModel
	// タイトル
	Title string `json:"title" gorm:"not null;check:title <> '';type:varchar(30)"`
	// 頻度ID
	FreqID uint `json:"freq_id"`
	// 面接フラグ
	InterviewFlg uint `json:"interview_flg"`
	// 開始時刻
	Start time.Time `json:"start" gorm:"not null"`
	// 終了時刻
	End time.Time `json:"end" gorm:"not null"`
	// 頻度(外部キー)
	CalendarFreqStatus CalendarFreqStatus `gorm:"foreignKey:freq_id;references:id"`
}

/*
	t_user_schedule_association
	ユーザー予定紐づけ
*/
type UserScheduleAssociation struct {
	// ユーザー予定ID
	UserScheduleID uint `json:"user_schedule_id" gorm:"primaryKey"`
	// ユーザーID
	UserID uint `json:"user_id" gorm:"primaryKey"`
	// ユーザー予定(外部キー)
	UserSchedule UserSchedule `gorm:"foreignKey:user_schedule_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

func (t User) TableName() string {
	return "t_user"
}
func (t UserGroup) TableName() string {
	return "t_user_group"
}
func (t UserGroupAssociation) TableName() string {
	return "t_user_group_association"
}
func (t UserSchedule) TableName() string {
	return "t_user_schedule"
}
func (t UserScheduleAssociation) TableName() string {
	return "t_user_schedule_association"
}

// ユーザーグループ登録 Request
type UserGroupRequest struct {
	UserGroup
	// 所属ユーザー
	Users string `json:"users"`
}

// ユーザー予定登録 Request
type UserScheduleCreateRequest struct {
	UserSchedule
	// ハッシュキー(ユーザー)
	UserHashKeys string `json:"user_hash_keys"`
}

// ユーザー予定 Request
type UserScheduleRequest struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// 応募者ハッシュキー
	ApplicantHashKey string `json:"applicant_hash_key"`
	// ハッシュキーリスト(ユーザー)
	UserHashKeys string `json:"user_hash_keys"`
	// 面接フラグ
	InterviewFlg uint `json:"interview_flg"`
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
	// 面接フラグ
	InterviewFlg uint `json:"interview_flg"`
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

// 予約時間
type ReserveTime struct {
	// 時間
	Time time.Time `json:"time"`
	// 予約可否
	IsReserve bool `json:"is_reserve"`
}

// 予約表
type ReserveTable struct {
	// 年月日
	Dates []time.Time `json:"date"`
	// 予約時間
	Options []ReserveTime `json:"options"`
}

// グループ毎の面接可能人数
type ReserveOfGroup struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// 面接可能人数
	Count uint `json:"count"`
}

type UserScheduleAssociationWithName struct {
	UserScheduleAssociation
	// 氏名
	Name string `json:"name"`
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
