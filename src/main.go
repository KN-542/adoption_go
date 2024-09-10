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
	teamRepository := repository.NewTeamRepository(db)
	scheduleRepository := repository.NewScheduleRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	companyRepository := repository.NewCompanyRepository(db)
	applicantRepository := repository.NewApplicantRepository(db, redis)

	// Validator
	commonValidator := validator.NewCommonValidator()
	userValidator := validator.NewUserValidator()
	teamValidator := validator.NewTeamValidator()
	scheduleValidator := validator.NewScheduleValidator()
	applicantValidator := validator.NewApplicantValidator()
	loginValidator := validator.NewLoginValidator()
	companyValidator := validator.NewCompanyValidator()
	roleValidator := validator.NewRoleValidator()
	manuscriptValidator := validator.NewManuscriptValidator()

	// Service
	commonService := service.NewCommonService(
		masterRepository,
		roleRepository,
		userRepository,
		teamRepository,
		commonValidator,
		redisRepository,
	)
	companyService := service.NewCompanyService(
		companyRepository,
		masterRepository,
		roleRepository,
		userRepository,
		teamRepository,
		companyValidator,
		dbRepository,
	)
	applicantService := service.NewApplicantService(
		applicantRepository,
		userRepository,
		teamRepository,
		scheduleRepository,
		masterRepository,
		awsRepository,
		googleRepository,
		redisRepository,
		applicantValidator,
		dbRepository,
		outerRepository,
	)
	loginService := service.NewLoginService(
		userRepository,
		teamRepository,
		applicantRepository,
		redisRepository,
		loginValidator,
		userValidator,
		dbRepository,
	)
	userService := service.NewUserService(
		userRepository,
		teamRepository,
		scheduleRepository,
		roleRepository,
		applicantRepository,
		manuscriptRepository,
		masterRepository,
		userValidator,
		teamValidator,
		dbRepository,
		outerRepository,
		redisRepository,
	)
	teamService := service.NewTeamService(
		dbRepository,
		redisRepository,
		userRepository,
		teamRepository,
		scheduleRepository,
		applicantRepository,
		roleRepository,
		manuscriptRepository,
		masterRepository,
		teamValidator,
		outerRepository,
	)
	scheduleService := service.NewScheduleService(
		dbRepository,
		redisRepository,
		userRepository,
		teamRepository,
		scheduleRepository,
		applicantRepository,
		roleRepository,
		manuscriptRepository,
		masterRepository,
		scheduleValidator,
		outerRepository,
	)
	manuscriptService := service.NewManuscriptService(
		manuscriptRepository,
		masterRepository,
		userRepository,
		teamRepository,
		dbRepository,
		redisRepository,
		manuscriptValidator,
	)
	roleService := service.NewRoleService(roleRepository, redisRepository, roleValidator)

	// Controller
	commonController := controller.NewCommonController(commonService, loginService)
	companyController := controller.NewCompanyController(companyService, loginService, roleService)
	loginController := controller.NewLoginController(loginService)
	applicantController := controller.NewApplicantController(applicantService, userService, loginService, roleService)
	userController := controller.NewUserController(userService, applicantService, loginService, roleService)
	teamController := controller.NewTeamController(teamService, applicantService, loginService, roleService)
	scheduleController := controller.NewScheduleController(scheduleService, applicantService, loginService, roleService)
	roleController := controller.NewRoleController(roleService, loginService)
	manuscriptController := controller.NewManuscriptController(manuscriptService, loginService, roleService)

	e := router.NewRouter(
		commonController,
		loginController,
		userController,
		teamController,
		scheduleController,
		companyController,
		applicantController,
		manuscriptController,
		roleController,
	)
	e.Logger.Fatal(e.Start(":8080"))
}
