package model

/*
t_role
ロール
*/
type CustomRole struct {
	AbstractTransactionModel
	AbstractTransactionFlgModel
}

/*
	t_role_association
	ロール紐づけ
*/
type RoleAssociation struct {
	// ロールID
	RoleID uint `json:"role_id" gorm:"primaryKey"`
	// ユーザーID
	UserID uint `json:"user_id" gorm:"primaryKey"`
	// ロール(外部キー)
	Role CustomRole `gorm:"foreignKey:role_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

func (t CustomRole) TableName() string {
	return "t_role"
}
func (t RoleAssociation) TableName() string {
	return "t_role_association"
}
