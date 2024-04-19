package model

/*
t_role
ロール
*/
type CustomRole struct {
	AbstractTransactionModel
	AbstractTransactionFlgModel
	// 付与ロール
	Roles string `json:"roles" gorm:"type:text"`
}

func (t CustomRole) TableName() string {
	return "t_role"
}
