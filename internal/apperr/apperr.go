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
	CodeRecordNotFound      Code = "record_not_found"
	CodeDuplicateEmail      Code = "duplicate_email"
	CodeInsufficientBalance Code = "insufficient_balance"
	CodeSameAccountTransfer Code = "same_account_transfer"
	CodeEditConflict        Code = "edit_conflict"
	CodeInvalidToken        Code = "invalid_token"
	CodeInvalidRequest      Code = "invalid_request"
	CodeInvalidCredentials  Code = "invalid_credentials"
	CodeUnauthenticated     Code = "unauthenticated"
	CodeAccountInactive     Code = "account_inactive"
	CodePermissionDenied    Code = "permission_denied"
	CodeRateLimited         Code = "rate_limited"
	CodeInternal            Code = "internal"
)

type Error struct {
	code   Code
	kind   Kind
	fields map[string]string
	cause  error
}

func New(code Code) *Error
func Field(code Code, field, msg string) *Error
func Wrap(code Code, cause error) *Error
func (e *Error) Code() Code
func (e *Error) Kind() Kind
func (e *Error) Fields() map[string]string
func (e *Error) Unwrap() error
func (e *Error) Error() string
