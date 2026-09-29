package api

import (
	"errors"
	"net/http"

	"github.com/nongpal/Palge-Backend/internal/apperr"
)

type descriptor struct {
	message string
}

type errorBody struct {
	Code      apperr.Code       `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

func (app *Application) writeError(w http.ResponseWriter, r *http.Request, err error) {
	aerr := &apperr.Error{}
	if !errors.As(err, &aerr) {
		aerr = apperr.New(apperr.CodeInternal)
	}

	d, ok := catalog[aerr.Code()]
	if !ok {
		app.logger.ErrorContext(r.Context(), "unmapped error code",
			"code", aerr.Code(), "request_id", app.contextGetRequestID(r.Context()))
		aerr = apperr.New(apperr.CodeInternal)
		d = catalog[apperr.CodeInternal]
	}

	requestID := app.contextGetRequestID(r.Context())
	if aerr.Kind() == apperr.KindInternal {
		app.logger.ErrorContext(r.Context(), "request failed",
			"error", err, "code", aerr.Code(),
			"request_id", requestID,
		)
	}

	status := kindStatus[aerr.Kind()]
	if aerr.Kind() == apperr.KindAuthentication {
		w.Header().Set("WWW-Authenticate", "Bearer")
	}

	if status == 0 {
		status = http.StatusInternalServerError
	}

	message := d.message
	if aerr.Kind() != apperr.KindInternal {
		if detail := aerr.Detail(); detail != "" {
			message = detail
		}
	}

	body := errorBody{
		Code:      aerr.Code(),
		Message:   message,
		Fields:    aerr.Fields(),
		RequestID: requestID,
	}

	if werr := app.writeJSON(w, status, envelope{"error": body}, nil); werr != nil {
		app.logger.ErrorContext(r.Context(), "failed to write error response", "error", werr, "request_id", requestID)
	}
}

var Codes = []apperr.Code{
	apperr.CodeRecordNotFound,
	apperr.CodeDuplicateEmail,
	apperr.CodeInsufficientBalance,
	apperr.CodeSameAccountTransfer,
	apperr.CodeEditConflict,
	apperr.CodeInvalidToken,
	apperr.CodeInvalidRequest,
	apperr.CodeValidationFailed,
	apperr.CodeInvalidCredentials,
	apperr.CodeUnauthenticated,
	apperr.CodeAccountInactive,
	apperr.CodePermissionDenied,
	apperr.CodeRateLimited,
	apperr.CodeInternal,
}

var catalog = map[apperr.Code]descriptor{
	apperr.CodeRecordNotFound: {
		message: "the requested resource could not be found",
	},

	apperr.CodeDuplicateEmail: {
		message: "a user with this email address already exists",
	},

	apperr.CodeInsufficientBalance: {
		message: "the account has insufficient balance for this operation",
	},

	apperr.CodeSameAccountTransfer: {
		message: "sender and receiver must be different accounts",
	},

	apperr.CodeEditConflict: {
		message: "unable to update the record due to an edit conflict, please try again",
	},

	apperr.CodeInvalidToken: {
		message: "invalid or expired token",
	},

	apperr.CodeInvalidRequest: {
		message: "invalid request",
	},

	apperr.CodeValidationFailed: {
		message: "the request body contains invalid fields",
	},

	apperr.CodeInvalidCredentials: {
		message: "invalid authentication credentials",
	},

	apperr.CodeUnauthenticated: {
		message: "you must be authenticated to access this resource",
	},

	apperr.CodeAccountInactive: {
		message: "your user account must be activated to access this resource",
	},

	apperr.CodePermissionDenied: {
		message: "your account does not have the necessary permissions to access this resource",
	},

	apperr.CodeRateLimited: {
		message: "rate limit exceeded",
	},

	apperr.CodeInternal: {
		message: "the server encountered a problem and could not process your request",
	},
}

var kindStatus = map[apperr.Kind]int{
	apperr.KindValidation:       http.StatusUnprocessableEntity,
	apperr.KindNotFound:         http.StatusNotFound,
	apperr.KindConflict:         http.StatusConflict,
	apperr.KindMalformedRequest: http.StatusBadRequest,
	apperr.KindAuthentication:   http.StatusUnauthorized,
	apperr.KindAuthorization:    http.StatusForbidden,
	apperr.KindRateLimit:        http.StatusTooManyRequests,
	apperr.KindInternal:         http.StatusInternalServerError,
}
