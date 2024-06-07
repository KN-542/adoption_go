package ddl

/*
t_notice
通知
*/
type Notice struct {
	AbstractTransactionModel
	// 種別
	Type uint `json:"type"`
	// 通知元ユーザーID
	FromUserID uint64 `json:"from_user_id"`
	// 通知先ユーザーID
	ToUserID uint64 `json:"to_user_id"`
	// 通知種別(外部キー)
	NoticeType NoticeType `gorm:"foreignKey:type;references:id"`
	// 通知元ユーザー(外部キー)
	FromUser User `gorm:"foreignKey:from_user_id;references:id"`
	// 通知先ユーザー(外部キー)
	ToUser User `gorm:"foreignKey:to_user_id;references:id"`
}

func (t Notice) TableName() string {
	return "t_notice"
}

// 共通Entity ※削除予定
type CommonModel struct {
	// ID
	ID uint64 `json:"id"`
	// ハッシュキー
	HashKey uint `json:"hash_key"`
}
