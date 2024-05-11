package ddl

type LoginStatusResponse struct {
	MFA            int8 `json:"mfa"`
	PasswordChange int8 `json:"password_change"`
}
