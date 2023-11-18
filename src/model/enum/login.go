package enum

type LoginStatus int8
type PasswordChangeFlg int8

const (
	MFA_UNAUTHENTICATED LoginStatus = 0
	MFA_AUTHENTICATED   LoginStatus = 1
	PASSWORD_CHANGE     LoginStatus = 2
)
const (
	PASSWORD_CHANGE_UNREQUIRED PasswordChangeFlg = 0
	PASSWORD_CHANGE_REQUIRED   PasswordChangeFlg = 1
)
