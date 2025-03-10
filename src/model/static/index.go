package static

// Response Body コード
const (
	/*
		共通
	*/
	CODE_BAD_REQUEST uint = 101

	/*
		login
	*/
	// ログイン認証
	CODE_LOGIN_AUTH uint = 1
	// MFA作成
	CODE_LOGIN_REQUIRED uint = 1
	// MFA
	CODE_INVALID_CODE uint = 1
	CODE_EXPIRED      uint = 2
	// パスワード変更
	CODE_INIT_PASSWORD_INCORRECT uint = 1
	// チーム未検証
	CODE_CONFIRM_TEAM_NOT_EXIST uint = 1
	// 応募者チェック
	CODE_CHECK_APPLICANT_TEST_FINISHED          uint = 11
	CODE_CHECK_APPLICANT_CANNOT_UPDATE_SCHEDULE uint = 12

	/*
		company
	*/
	// 登録
	CODE_COMPANY_NAME_DUPL  uint = 1
	CODE_COMPANY_EMAIL_DUPL uint = 2

	/*
		ユーザー
	*/
	// 登録
	CODE_USER_EMAIL_DUPL uint = 1
	// チーム登録
	CODE_TEAM_USER_NOT_FOUND uint = 1
	// チーム削除
	CODE_TEAM_USER_CANNOT_DELETE_APPLICANT  uint = 1
	CODE_TEAM_USER_CANNOT_DELETE_SCHEDULE   uint = 2
	CODE_TEAM_USER_CANNOT_DELETE_MANUSCRIPT uint = 3
	CODE_TEAM_USER_CANNOT_DELETE_TEAM       uint = 4
	// ユーザー削除
	CODE_USER_CANNOT_DELETE_APPLICANT uint = 1
	CODE_USER_CANNOT_DELETE_SCHEDULE  uint = 2
	CODE_USER_CANNOT_DELETE_SELF      uint = 3

	/*
		応募者
	*/
	// 面接希望日登録
	CODE_APPLICANT_CANNOT_ASSIGN_USER uint = 1
	// 面接官割り振り
	CODE_APPLICANT_SCHEDULE_DOES_NOT_EXIST uint = 1
	CODE_APPLICANT_SHORTAGE_USER_MIN       uint = 2

	/*
		原稿
	*/
	// 原稿登録
	CODE_MANUSCRIPT_DUPLICATE_CONTENT uint = 1
	// 原稿削除
	CODE_MANUSCRIPT_CANNOT_DELETE_APPLICANT uint = 1
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
