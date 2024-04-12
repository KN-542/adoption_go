package model

import "time"

/*
t_operation_log
操作ログ
*/
type OperationLog struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// イベントID
	EventID uint `json:"event_id"`
	// 対象ユーザーハッシュキー
	UserHashKey string `json:"user_hash_key" gorm:"check:user_hash_key <> '';type:text"`
	// ログ
	Log string `json:"log" gorm:"type:text"`
	// 企業ID
	CompanyID uint `json:"company_id" gorm:"index"`
	// 登録日時
	CreatedAt time.Time `json:"created_at"`
	// イベント
	Event OperationLogEvent `gorm:"foreignKey:event_id;references:id"`
	// 企業(外部キー)
	Company Company `gorm:"foreignKey:company_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_hash_key;references:hash_key"`
}
