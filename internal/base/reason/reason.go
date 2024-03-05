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
	PasswordNotMatch          = "error.passport.password_not_match"
	InvalidEmail              = "error.passport.invalid_email"
	CannotDeleteCurrentDevice = "error.passport.cannot_delete_current_device"
	DeviceNotExist            = "error.passport.device_not_exist"

	UserInvalidName     = "error.user.invalid_name"
	UserInvalidPassword = "error.user.invalid_password"
	UserNotExist        = "error.user.not_exist"
	UserDuplicate       = "error.user.duplicate"

	EmailNotExist  = "errors.email.not_exist"
	EmailDuplicate = "errors.email.duplicate"

	CaptchaNotExist       = "errors.capcha.not_exist"
	CaptchaMismatch       = "errors.capcha.mismatch"
	CaptchaTypeNotSupport = "errors.capcha.type_not_support"
	CaptchaReachSendLimit = "errors.capcha.reach_send_limit"
)
