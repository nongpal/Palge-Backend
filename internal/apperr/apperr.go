package apperr

import "fmt"

type Code string
type Kind int

const (
	KindValidation Kind = iota
	KindNotFound
	KindConflict
	KindMalformedRequest
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

var codeKind = map[Code]Kind{
	CodeRecordNotFound:      KindNotFound,
	CodeDuplicateEmail:      KindValidation,
	CodeInsufficientBalance: KindValidation,
	CodeSameAccountTransfer: KindValidation,
	CodeEditConflict:        KindConflict,
	CodeInvalidToken:        KindAuthentication,
	CodeInvalidRequest:      KindMalformedRequest,
	CodeInvalidCredentials:  KindAuthentication,
	CodeUnauthenticated:     KindAuthentication,
	CodeAccountInactive:     KindAuthorization,
	CodePermissionDenied:    KindAuthorization,
	CodeRateLimited:         KindRateLimit,
	CodeInternal:            KindInternal,
}

type Error struct {
	code   Code
	kind   Kind
	fields map[string]string
	cause  error
}

func New(code Code) *Error {
	return &Error{
		code: code,
		kind: codeKind[code],
	}
}

func Field(code Code, field, msg string) *Error {
	e := New(code)
	e.fields = map[string]string{field: msg}
	return e
}

func Wrap(code Code, cause error) *Error {
	e := New(code)
	e.cause = cause
	return e
}

func (e *Error) Code() Code {
	return e.code
}

func (e *Error) Kind() Kind {
	return e.kind
}

func (e *Error) Fields() map[string]string {
	return e.fields
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.code, e.cause)
	}
	return string(e.code)
}
