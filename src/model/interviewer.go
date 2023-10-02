package model

import "time"

type Interviewer struct {
    ID   string `json:"id" gorm:"primaryKey;type:varchar(255)"`
    Name string `json:"name" gorm:"notNull;type:varchar(255)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InterviewerResponse struct {
    ID   string `json:"id" gorm:"primaryKey;type:varchar(255)"`
    Name string `json:"name" gorm:"notNull;type:varchar(255)"`
}

func (t Interviewer) TableName() string {
    return "t_interviewer"
}