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

	dbRepository := repository.NewDBRepository(db)

	redis := infra.NewRedis()
	redisRepository := repository.NewRedisRepository(redis)

	awsRepository := repository.NewAWSRepository()

	masterRepository := repository.NewMasterRepository(db)

	userRepository := repository.NewUserRepository(db)
	userValidate := validator.NewUserValidator()

	applicantRepository := repository.NewApplicantRepository(db, redis)
	applicantValidator := validator.NewApplicantValidator()
	applicantService := service.NewApplicantService(applicantRepository, masterRepository, awsRepository, redisRepository, applicantValidator, dbRepository)
	applicantController := controller.NewApplicantController(applicantService)

	loginService := service.NewLoginService(userRepository, applicantRepository, redisRepository, userValidate, dbRepository)
	loginController := controller.NewLoginController(loginService)

	userService := service.NewUserService(userRepository, masterRepository, userValidate, dbRepository)
	userController := controller.NewUserController(userService)

	e := router.NewRouter(
		loginController,
		userController,
		applicantController,
	)
	e.Logger.Fatal(e.Start(":8080"))
}
