package main

import (
	"api/src/controller"
	"api/src/infra"
	"api/src/repository"
	"api/src/router"
	"api/src/service"
	"api/src/validator"
)

func main() {
	db := infra.NewDB()

	redis := infra.NewRedis()
	redisRepository := repository.NewRedisRepository(redis)

	masterRepository := repository.NewMasterRepository(db)

	userRepository := repository.NewUserRepository(db)
	userValidate := validator.NewUserValidator()

	loginService := service.NewLoginService(userRepository, redisRepository, userValidate)
	loginController := controller.NewLoginController(loginService)

	userService := service.NewUserService(userRepository, masterRepository, userValidate)
	userController := controller.NewUserController(userService)

	applicantRepository := repository.NewApplicantRepository(db, redis)
	applicantService := service.NewApplicantService(applicantRepository, masterRepository)
	applicantController := controller.NewApplicantController(applicantService)

	e := router.NewRouter(
		loginController,
		userController,
		applicantController,
	)
	e.Logger.Fatal(e.Start(":8080"))
}
