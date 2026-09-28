package apperr

type Code string
type Kind int

const (
	KindValidation Kind = iota
	KindNotFound
	KindConflict
	KindAuthentication
	KindAuthorization
	KindRateLimit
	KindInternal
)

const (
	CodeRecordNotFound Code = "record_not_found"
	CodeDuplicateEmail Code = "duplicate_email"
	CodeInsufficientBalance Code = "insufficient_balance"
	CodeSameAccountTransfer Code = "same_account_transfer"
	CodeEditConflict Code = "edit_conflict"
	CodeInvalidToken Code = "invalid_token"
	CodeInvalidRequest Code = "invalid_request"
	CodeInvalidCredentials Code = "invalid_credentials"
	CodeUnauthenticated Code = "unauthenticated"
	CodeAccountInactive Code = "account_inactive"
	CodePermissionDenied Code = "permission_denied"
	CodeRateLimited Code = "rate_limited"
	CodeInternal Code = "internal"
)