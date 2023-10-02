package service

import (
	"api/src/model"
	"api/src/repository"
	"api/src/validator"
	"time"

	"github.com/google/uuid"
)

type IInterviewerService interface {
	RegisterInterviewer(i model.Interviewer) error
	Interviewers() ([]model.InterviewerResponse, error)
}

type interviewerService struct {
	r repository.IInterviewerRepository
	v validator.IInterviewerValidator
}

func NewInterviewerService(r repository.IInterviewerRepository, v validator.IInterviewerValidator) IInterviewerService {
	return &interviewerService{r, v}
}

func (s *interviewerService) RegisterInterviewer(i model.Interviewer) error {
	// バリデーション
	if err := s.v.InterviewerValidate(i); err != nil {
		return err
	}

	// ID生成
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	now := time.Now().Format("20060102150405")
	id := uuid.String() + now

	interviewer := model.Interviewer{ID: id, Name: i.Name, CreatedAt: i.CreatedAt, UpdatedAt: i.UpdatedAt}
	if err := s.r.Insert(&interviewer); err != nil {
		return err
	}

	return nil
}

func (s *interviewerService) Interviewers() ([]model.InterviewerResponse, error) {
	resList, err := s.r.List()
	if err != nil {
		return []model.InterviewerResponse{}, err
	}
	return resList, nil
}