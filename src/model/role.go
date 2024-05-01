package model

/*
t_role
ロール
*/
type CustomRole struct {
	AbstractTransactionModel
	AbstractTransactionFlgModel
	// ロール名
	Name string `json:"name" gorm:"unique;not null;check:name <> '';type:varchar(75);index"`
}

/*
	t_role_association
	付与ロール
*/
type RoleAssociation struct {
	// ロールID
	RoleID uint `json:"role_id" gorm:"primaryKey"`
	// マスターロールID
	MasterRoleID uint `json:"master_role_id" gorm:"primaryKey"`
	// ロール(外部キー)
	CustomRole CustomRole `gorm:"foreignKey:role_id;references:id"`
	// ロールマスタ(外部キー)
	Role Role `gorm:"foreignKey:master_role_id;references:id"`
}

func (t CustomRole) TableName() string {
	return "t_role"
}
func (t RoleAssociation) TableName() string {
	return "t_role_association"
}
