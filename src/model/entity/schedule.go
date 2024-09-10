package entity

import "api/src/model/ddl"

// Schedule
type Schedule struct {
	ddl.Schedule
	// 頻度名
	FreqName string `json:"freq_name"`
	// 該当ユーザー
	Users []*User `json:"users" gorm:"many2many:t_schedule_association;foreignKey:id;joinForeignKey:schedule_id;References:id;joinReferences:user_id"`
}

// Schedule
type Schedule2 struct {
	ddl.Schedule
}

// ScheduleAssociation
type ScheduleAssociation struct {
	ddl.ScheduleAssociation
}
