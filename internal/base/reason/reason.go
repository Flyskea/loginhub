package reason

const (
	// Success .
	Success = "base.success"
	// UnknownError unknown error
	UnknownError = "base.unknown"
	// RequestFormatError request format error
	RequestFormatError = "base.request_format_error"
	// UnauthorizedError unauthorized error
	UnauthorizedError = "base.unauthorized_error"
	// DatabaseError database error
	DatabaseError = "base.database_error"
	// ForbiddenError forbidden error
	ForbiddenError = "base.forbidden_error"
	// RateLimitRequestError rate limit request error
	RateLimitRequestError = "base.rate_limit_request_error"
	// OjbectNotFoundError object not found
	ObjectNotFoundError = "base.object_not_found"
)

const (
	AccountOrPasswordWrong    = "error.passport.account_or_password_incorrect"
	RegisterTypeNotSupport    = "error.passport.register_type_not_support"
	LoginTypeNotSupport       = "error.passport.login_type_not_support"
	PasswordMismatch          = "error.passport.password_mismatch"
	InvalidEmail              = "error.passport.invalid_email"
	CannotDeleteCurrentDevice = "error.passport.cannot_delete_current_device"
	DeviceNotExist            = "error.passport.device_not_exist"

	UserInvalidName     = "error.user.invalid_name"
	UserInvalidPassword = "error.user.invalid_password"
	UserNotExist        = "error.user.not_exist"
	UserDuplicate       = "error.user.duplicate"

	OAuth2RedirectURLEmpty       = "error.oauth2.redirect_url_empty"
	OAuth2ClientIDEmpty          = "error.oauth2.client_id_empty"
	OAuth2ClientSecretEmpty      = "error.oauth2.client_secret_empty"
	OAuth2ProviderNotImplemented = "error.oauth2.provider_not_implemented"
	OAuth2ProviderDuplicate      = "error.oauth2.provider_duplicate"
	OAuth2ProviderNotExist       = "error.oauth2.provider_not_exist"
	OAuth2StateEmpty             = "error.oauth2.state_empty"
	OAuth2StateMismatch          = "error.oauth2.state_mismatch"
	OAuth2StateNotFound          = "error.oauth2.state_not_found"

	EmailNotExist  = "errors.email.not_exist"
	EmailDuplicate = "errors.email.duplicate"

	CaptchaNotExist       = "errors.capcha.not_exist"
	CaptchaMismatch       = "errors.capcha.mismatch"
	CaptchaTypeNotSupport = "errors.capcha.type_not_support"
	CaptchaReachSendLimit = "errors.capcha.reach_send_limit"
)
