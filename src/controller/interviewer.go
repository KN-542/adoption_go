package controller

import (
	"api/src/model"
	"api/src/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IInterviewerController interface {
	RegisterInterviewer(e echo.Context) error
	Interviewers(e echo.Context) error
}

type interviewerController struct {
	i service.IInterviewerService
}

func NewInterviewerController(i service.IInterviewerService) IInterviewerController {
	return &interviewerController{i}
}

func (c *interviewerController) RegisterInterviewer(e echo.Context) error {
	interviewer := model.Interviewer{}
	if err := e.Bind(&interviewer); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.i.RegisterInterviewer(interviewer); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, nil)
}

func (c *interviewerController) Interviewers(e echo.Context) error {
	interviewer := model.Interviewer{}
	if err := e.Bind(&interviewer); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := c.i.Interviewers()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, res)
}