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