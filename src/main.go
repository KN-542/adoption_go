package main

import (
	"api/src/controller"
	"api/src/infra/db"
	"api/src/infra/redis"
	"api/src/repository"
	"api/src/router"
	"api/src/service"
	"api/src/validator"
	"log"
)

func main() {
	db := db.NewDB()

	redis := redis.NewRedis()

	masterRepository := repository.NewMasterRepository(db)

	loginService := service.NewLoginService()
	loginController := controller.NewLoginController(loginService)

	userRepository := repository.NewUserRepository(db)
	userValidate := validator.NewUserValidator()
	userService := service.NewUserService(userRepository, masterRepository, userValidate)
	userController := controller.NewUserController(userService)

	applicantRepository := repository.NewApplicantRepository(db, redis)
	applicantService := service.NewApplicantService(applicantRepository, masterRepository)
	applicantController := controller.NewApplicantController(applicantService)

	log.Print(222)
	e := router.NewRouter(
		loginController,
		userController,
		applicantController,
	)
	log.Print(333)
	e.Logger.Fatal(e.Start(":8080"))
}
