package validator

// import (
// 	"api/src/model"

// 	validation "github.com/go-ozzo/ozzo-validation/v4"
// )

// type IInterviewerValidator interface {
// 	InterviewerValidate(i model.Interviewer) error
// }

// type interviewerValidator struct {}

// func NewInterviewerValidator() IInterviewerValidator {
// 	return &interviewerValidator{}
// }

// func (v *interviewerValidator) InterviewerValidate(i model.Interviewer) error {
// 	return validation.ValidateStruct(&i, validation.Field(&i.Name, validation.Required.Error("name is required")))
// }
