package static

// Response Body コード
const (
	/*
		common
	*/
	CODE_BAD_REQUEST uint = 101

	/*
		login
	*/
	// Login
	CODE_LOGIN_AUTH uint = 1
	// MFA create && Session confirm
	CODE_LOGIN_REQUIRED uint = 1
	// MFA
	CODE_INVALID_CODE uint = 1
	CODE_EXPIRED      uint = 2
	// PasswordChange
	CODE_INIT_PASSWORD_INCORRECT uint = 1

	/*
		user
	*/
	// create
	CODE_COMPANY_NAME_DUPL uint = 1

	/*
		user
	*/
	// create
	CODE_USER_EMAIL_DUPL uint = 1
)

// Response Body メッセージ
const (
	/*
		common
	*/
	MESSAGE_BAD_REQUEST             string = "bad request"
	MESSAGE_UNEXPECTED_COOKIE       string = "unexpected jwt token"
	MESSAGE_NOT_FOUND_LOGIN_SERVICE string = "controller does not have a valid ILoginService field"
)

// 言語
const (
	JA string = "ja"
)

var Langs = []string{JA}

// Oauth2.0
const (
	OAUTH_CODE          string = "code"
	OAUTH_CLIENT_ID     string = "client_id"
	OAUTH_CLIENT_SECRET string = "client_secret"
	OAUTH_GRANT_TYPE    string = "grant_type"
	OAUTH_REDIRECT_URI  string = "redirect_uri"
	OAUTH_ACCESS_TOKEN  string = "access_token"
	OAUTH_REFRESH_TOKEN string = "refresh_token"
)

// Redis Key
const (
	// ユーザー
	REDIS_USER_HASH_KEY   string = "hash_key"
	REDIS_USER_ROLE       string = "role_id"
	REDIS_USER_LOGIN_TYPE string = "login_type"
	REDIS_USER_COMPANY_ID string = "company_id"
	REDIS_USER_TEAM_ID    string = "team_id"
	// 応募者
	REDIS_APPLICANT_HASH_KEY  string = "applicant_hash_key"
	REDIS_CODE                string = "code"
	REDIS_S3_NAME             string = "s3_name"
	REDIS_OAUTH_REFRESH_TOKEN string = "oauth_refresh_token"
)