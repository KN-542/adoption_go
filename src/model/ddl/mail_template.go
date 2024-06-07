package ddl

import "time"

/*
t_mail_template
メールテンプレート
*/
type MailTemplate struct {
	AbstractTransactionModel
	// メールテンプレート名
	Title string `json:"title" gorm:"not null;unique;check:title <> '';type:varchar(50)"`
	// 件名
	Subject string `json:"subject" gorm:"not null;unique;check:subject <> '';type:text"`
	// テンプレート
	Template string `json:"template" gorm:"not null;unique;check:template <> '';type:text"`
	// 説明
	Desc string `json:"desc" gorm:"type:text"`
	AbstractTransactionFlgModel
}

/*
t_variable
変数
*/
type Variable struct {
	AbstractTransactionModel
	// 変数タイトル
	Title string `json:"title" gorm:"not null;unique;check:title <> '';type:varchar(25)"`
	// 変数格納Json名
	JsonName string `json:"json_name" gorm:"not null;unique;check:json_name <> '';type:text"`
	AbstractTransactionFlgModel
}

/*
t_mail_preview
メールプレビュー
*/
type MailPreview struct {
	// テンプレートID
	TemplateID uint64 `json:"template_id" gorm:"primaryKey"`
	// 変数ID
	VariableID uint64 `json:"variable_id" gorm:"primaryKey"`
	// ハッシュキー
	HashKey string `json:"hash_key" gorm:"not null;unique;check:hash_key <> '';type:text;index"`
	// メールプレビュー名
	Title string `json:"title" gorm:"not null;unique;check:title <> '';type:varchar(50)"`
	// 説明
	Desc string `json:"desc" gorm:"type:text"`
	// 企業ID
	CompanyID uint64 `json:"company_id" gorm:"index"`
	// 登録日時
	CreatedAt time.Time `json:"created_at"`
	// 更新日時
	UpdatedAt time.Time `json:"updated_at"`
	// メールテンプレート(外部キー)
	MailTemplate MailTemplate `gorm:"foreignKey:template_id;references:id"`
	// 変数(外部キー)
	Variable Variable `gorm:"foreignKey:variable_id;references:id"`
	// 企業(外部キー)
	Company Company `gorm:"foreignKey:company_id;references:id"`
	AbstractTransactionFlgModel
}

func (t MailTemplate) TableName() string {
	return "t_mail_template"
}
func (t Variable) TableName() string {
	return "t_variable"
}
func (t MailPreview) TableName() string {
	return "t_mail_preview"
}
