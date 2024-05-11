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
	// DB
	db := infra.NewDB()

	// Redis
	redis := infra.NewRedis()

	// Repository
	dbRepository := repository.NewDBRepository(db)
	redisRepository := repository.NewRedisRepository(redis)
	outerRepository := repository.NewOuterRepository()
	awsRepository := repository.NewAWSRepository()
	googleRepository := repository.NewGoogleRepository(redis)
	masterRepository := repository.NewMasterRepository(db)
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	companyRepository := repository.NewCompanyRepository(db)
	applicantRepository := repository.NewApplicantRepository(db, redis)

	// Validator
	commonValidator := validator.NewCommonValidator()
	userValidator := validator.NewUserValidator()
	applicantValidator := validator.NewApplicantValidator()
	loginValidator := validator.NewLoginValidator()
	companyValidator := validator.NewCompanyValidator()

	// Service
	commonService := service.NewCommonService(
		masterRepository,
		roleRepository,
		commonValidator,
		redisRepository,
	)
	companyService := service.NewCompanyService(
		companyRepository,
		masterRepository,
		roleRepository,
		userRepository,
		companyValidator,
		dbRepository,
	)
	applicantService := service.NewApplicantService(
		applicantRepository,
		userRepository,
		masterRepository,
		awsRepository,
		googleRepository,
		redisRepository,
		applicantValidator,
		dbRepository,
	)
	loginService := service.NewLoginService(userRepository, applicantRepository, redisRepository, loginValidator, userValidator, dbRepository)
	userService := service.NewUserService(userRepository, applicantRepository, masterRepository, userValidator, dbRepository, outerRepository)

	// Controller
	commonController := controller.NewCommonController(commonService, loginService)
	companyController := controller.NewCompanyController(companyService, loginService)
	loginController := controller.NewLoginController(loginService)
	applicantController := controller.NewApplicantController(applicantService, userService)
	userController := controller.NewUserController(userService)

	e := router.NewRouter(
		commonController,
		loginController,
		userController,
		companyController,
		applicantController,
	)
	e.Logger.Fatal(e.Start(":8080"))
}
