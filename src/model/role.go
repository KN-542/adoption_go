package model

import "time"

/*
t_role
ロール
*/
type CustomRole struct {
	// ID
	ID uint64 `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	// ハッシュキー
	HashKey string `json:"hash_key" gorm:"not null;unique;check:hash_key <> '';type:text"`
	// 付与ロール
	Roles string `json:"roles" gorm:"type:text"`
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

func (t CustomRole) TableName() string {
	return "t_role"
}
