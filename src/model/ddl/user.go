package ddl

/*
t_user
ユーザー
*/
type User struct {
	AbstractTransactionModel
	// 氏名
	Name string `json:"name" gorm:"not null;check:name <> '';type:varchar(75);index"`
	// メールアドレス
	Email string `json:"email" gorm:"not null;unique;type:varchar(100);check:email ~ '^[a-zA-Z0-9_+-]+(\\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$';index"`
	// パスワード(ハッシュ化)
	Password string `json:"password" gorm:"not null;check:password <> ''"`
	// 初回パスワード(ハッシュ化)
	InitPassword string `json:"init_password" gorm:"not null;check:init_password <> ''"`
	// ロールID
	RoleID uint64 `json:"role_id" gorm:"index"`
	// ユーザー種別
	UserType uint `json:"user_type" gorm:"index"`
	// ロール(外部キー)
	Role CustomRole `gorm:"foreignKey:role_id;references:id"`
	// ログイン種別(外部キー)
	LoginType LoginType `gorm:"foreignKey:user_type;references:id"`
}

/*
t_user_refresh_token_association
リフレッシュトークン紐づけ
*/
type UserRefreshTokenAssociation struct {
	// ユーザーID
	UserID uint64 `json:"user_id" gorm:"primaryKey"`
	// リフレッシュトークン
	RefreshToken string `json:"refresh_token" gorm:"type:text"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

func (t User) TableName() string {
	return "t_user"
}
func (t UserRefreshTokenAssociation) TableName() string {
	return "t_user_refresh_token_association"
}
