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
	manuscriptRepository := repository.NewManuscriptRepository(db)
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
	roleValidator := validator.NewRoleValidator()

	// Service
	commonService := service.NewCommonService(
		masterRepository,
		roleRepository,
		userRepository,
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
		outerRepository,
	)
	loginService := service.NewLoginService(userRepository, applicantRepository, redisRepository, loginValidator, userValidator, dbRepository)
	userService := service.NewUserService(
		userRepository,
		roleRepository,
		applicantRepository,
		manuscriptRepository,
		masterRepository,
		userValidator,
		dbRepository,
		outerRepository,
		redisRepository,
	)
	roleService := service.NewRoleService(roleRepository, redisRepository, roleValidator)

	// Controller
	commonController := controller.NewCommonController(commonService, loginService)
	companyController := controller.NewCompanyController(companyService, loginService, roleService)
	loginController := controller.NewLoginController(loginService)
	applicantController := controller.NewApplicantController(applicantService, userService, loginService, roleService)
	userController := controller.NewUserController(userService, applicantService, loginService, roleService)
	roleController := controller.NewRoleController(roleService, loginService)

	e := router.NewRouter(
		commonController,
		loginController,
		userController,
		companyController,
		applicantController,
		roleController,
	)
	e.Logger.Fatal(e.Start(":8080"))
}
