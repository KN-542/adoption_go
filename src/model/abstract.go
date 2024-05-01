package model

import "time"

type AbstractMasterModel struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
}

type AbstractTransactionModel struct {
	// ID
	ID uint64 `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	// ハッシュキー
	HashKey string `json:"hash_key" gorm:"not null;unique;check:hash_key <> '';type:text;index"`
	// 企業ID
	CompanyID uint `json:"company_id" gorm:"index"`
	// 登録日時
	CreatedAt time.Time `json:"created_at"`
	// 更新日時
	UpdatedAt time.Time `json:"updated_at"`
	// 企業(外部キー)
	Company Company `gorm:"foreignKey:company_id;references:id"`
}

type AbstractTransactionFlgModel struct {
	// 編集保護
	EditFlg uint `json:"edit_flg"`
	// 削除保護
	DeleteFlg uint `json:"delete_flg"`
}
