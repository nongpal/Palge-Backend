package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nongpal/Palge-Backend/internal/data"
)

type errorKind int

const (
	errorKindValidation errorKind = iota
	errorKindNotFound
	errorKindConflict
	errorKindAuthentication
	errorKindAuthorization
	errorKindRateLimit
	errorKindInternal
)

type HTTPError struct {
	status  int
	message any
	kind    errorKind
}

func classifyError(err error) HTTPError {
	var status int
	var message any
	var kind errorKind

	switch {
	case errors.Is(err, data.ErrRecordNotFound):
		status = http.StatusNotFound
		message = "the requested resource could not be found"
		kind = errorKindNotFound

	case errors.Is(err, data.ErrDuplicateEmail):
		status = http.StatusUnprocessableEntity
		message = map[string]string{
			"email": "a user with this email address already exists",
		}
		kind = errorKindValidation

	case errors.Is(err, data.ErrEditConflict):
		status = http.StatusConflict
		message = "unable to update the record due to an edit conflict, please try again"
		kind = errorKindConflict

	case errors.Is(err, data.ErrInsufficientBalance):
		status = http.StatusUnprocessableEntity
		message = "insufficient account balance"
		kind = errorKindValidation

	default:
		status = http.StatusInternalServerError
		message = "the server encountered a problem and could not process your request"
		kind = errorKindInternal

	}

	return HTTPError{
		status:  status,
		message: message,
		kind:    kind,
	}
}

func (app *Application) applicationErrorResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	herr := classifyError(err)

	if herr.kind == errorKindInternal {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.errorResponse(w, r, herr.status, herr.message)
}

func (app *Application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	envlp := envelope{"error": message}

	err := app.writeJSON(w, status, envlp, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	envlp := envelope{"error": message, "request_id": app.contextGetRequestID(r.Context())}

	err = app.writeJSON(w, http.StatusInternalServerError, envlp, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *Application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *Application) invalidCredentialResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *Application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *Application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
func (app *Application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

func (app *Application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := `your user account doesn't have the ncessary permissions to access this resource`
	app.errorResponse(w, r, http.StatusForbidden, message)
}

func (app *Application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	envlp := envelope{
		"error":      "rate limit exceeded",
		"request_id": app.contextGetRequestID(r.Context()),
	}
	err := app.writeJSON(w, http.StatusTooManyRequests, envlp, nil)
	if err != nil {
		app.logError(r, err)
	}
}
