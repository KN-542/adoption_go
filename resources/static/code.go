package static

// code
const (
	/*
		common
	*/
	CODE_BAD_REQUEST int8 = 101

	/*
		login
	*/
	// Login
	CODE_LOGIN_AUTH int8 = 1
	// MFA create && Session confirm
	CODE_LOGIN_REQUIRED int8 = 1
	// MFA
	CODE_INVALID_CODE int8 = 1
	CODE_EXPIRED      int8 = 2
	// PasswordChange
	CODE_INIT_PASSWORD_INCORRECT int8 = 1

	/*
		user
	*/
	// create
	CODE_COMPANY_NAME_DUPL int8 = 1

	/*
		user
	*/
	// create
	CODE_USER_EMAIL_DUPL int8 = 1
)

// message
const (
	/*
		common
	*/
	MESSAGE_BAD_REQUEST             string = "bad request"
	MESSAGE_UNEXPECTED_COOKIE       string = "unexpected jwt token"
	MESSAGE_NOT_FOUND_LOGIN_SERVICE string = "controller does not have a valid ILoginService field"
)
