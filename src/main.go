package main

import (
	"api/src/controller"
	"api/src/infra/db"
	"api/src/infra/redis"
	"api/src/repository"
	"api/src/router"
	"api/src/service"
	"api/src/validator"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := db.NewDB()

	redis := redis.NewRedis()

	masterRepository := repository.NewMasterRepository(db)

	loginController := controller.NewLoginController()

	interviewerValidator := validator.NewInterviewerValidator()
	interviewerRepository := repository.NewInterviewerRepository(db)
	interviewerService := service.NewInterviewerService(interviewerRepository, interviewerValidator)
	interviewerController := controller.NewInterviewerController(interviewerService)

	applicantRepository := repository.NewApplicantRepository(db, redis)
	applicantService := service.NewApplicantService(applicantRepository, masterRepository)
	applicantController := controller.NewApplicantController(applicantService)

	e := router.NewRouter(loginController, interviewerController, applicantController)

	// CORSミドルウェアの設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("FE_CSR_URL"), os.Getenv("FE_SSR_URL")},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.Logger.Fatal(e.Start(":8080"))
}
