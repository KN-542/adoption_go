package model

import "time"

/*
t_mail_template
メールテンプレート
*/
type MailTemplate struct {
	// ID
	ID uint64 `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	// ハッシュキー
	HashKey string `json:"hash_key" gorm:"not null;unique;check:hash_key <> '';type:text"`
	// 件名
	Subject string `json:"subject" gorm:"not null;unique;check:subject <> '';type:text"`
	// テンプレート
	Template string `json:"template" gorm:"not null;unique;check:template <> '';type:text"`
	// 説明
	Desc string `json:"desc" gorm:"type:text"`
	// 利用方針
	Method uint `json:"method" gorm:"index"`
	// 変数キー名
	Keys string `json:"keys" gorm:"type:text"`
	// 企業ID
	CompanyID uint `json:"company_id" gorm:"index"`
	// 編集可能フラグ
	EditFlg uint `json:"edit_flg"`
	// 削除可能フラグ
	DeleteFlg uint `json:"delete_flg"`
	// 登録日時
	CreatedAt time.Time `json:"created_at"`
	// 更新日時
	UpdatedAt time.Time `json:"updated_at"`
	// 企業(外部キー)
	Company Company `gorm:"foreignKey:company_id;references:id"`
}

/*
t_variable
変数
*/
type Variable struct {
	// ID
	ID uint64 `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	// ハッシュキー
	HashKey string `json:"hash_key" gorm:"not null;unique;check:hash_key <> '';type:text"`
	// テンプレート
	Template string `json:"template" gorm:"type:text"`
	// 企業ID
	CompanyID uint `json:"company_id" gorm:"index"`
	// 編集可能フラグ
	EditFlg uint `json:"edit_flg"`
	// 削除可能フラグ
	DeleteFlg uint `json:"delete_flg"`
	// 登録日時
	CreatedAt time.Time `json:"created_at"`
	// 更新日時
	UpdatedAt time.Time `json:"updated_at"`
	// 企業(外部キー)
	Company Company `gorm:"foreignKey:company_id;references:id"`
}

func (t MailTemplate) TableName() string {
	return "t_mail_template"
}
func (t Variable) TableName() string {
	return "t_variable"
}
