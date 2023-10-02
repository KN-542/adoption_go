package repository

import (
	"api/src/model"

	"gorm.io/gorm"
)

type IInterviewerRepository interface {
	List() ([]model.InterviewerResponse, error)
	Insert(i *model.Interviewer) error
}

type interviewerRepository struct {
	db *gorm.DB
}

func NewInterviewerRepository(db *gorm.DB) IInterviewerRepository {
	return &interviewerRepository{db}
}

func (r *interviewerRepository) List() ([]model.InterviewerResponse, error) {
	var list []model.Interviewer
	var resList []model.InterviewerResponse

	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}

	for _, obj := range list {
		resList = append(resList, model.InterviewerResponse{ID: obj.ID, Name: obj.Name})
	}
	return resList, nil
}

func (r *interviewerRepository) Insert(i *model.Interviewer) error {
	if err := r.db.Create(i).Error; err != nil {
		return err
	}
	return nil
}