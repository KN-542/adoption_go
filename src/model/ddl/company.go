package ddl

import "time"

/*
t_company
企業
*/
type Company struct {
	// ID
	ID uint64 `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	// ハッシュキー
	HashKey string `json:"hash_key" gorm:"unique;not null;unique;check:hash_key <> '';type:text;index"`
	// 企業名
	Name string `json:"name" gorm:"not null;unique;check:name <> '';type:varchar(30)"`
	// ロゴファイル名
	Logo string `json:"logo" gorm:"not null;type:varchar(30)"`
	// 削除フラグ
	DeleteFlg uint `json:"delete_flg"`
	// 登録日時
	CreatedAt time.Time `json:"created_at"`
	// 更新日時
	UpdatedAt time.Time `json:"updated_at"`
}

func (t Company) TableName() string {
	return "t_company"
}
