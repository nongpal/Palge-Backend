package api

import (
	"net/http"
	"time"

	"github.com/nongpal/Palge-Backend/internal/apperr"
	"github.com/nongpal/Palge-Backend/internal/data"
	"github.com/nongpal/Palge-Backend/internal/validator"
)

func (app *Application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.writeError(w, r, apperr.Wrap(apperr.CodeInvalidRequest, err))
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.writeError(w, r, apperr.WithFields(apperr.CodeValidationFailed, v.Errors))
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case apperr.Is(err, apperr.CodeRecordNotFound):
			app.recordAuditLog(r, &data.AuditLog{
				UserID:     nil,
				Action:     "login",
				Resource:   "user",
				ResourceID: nil,
				Result:     "failed",
			})

			app.writeError(w, r, apperr.New(apperr.CodeInvalidCredentials))
		default:
			app.writeError(w, r, err)
		}
		return
	}

	match, err := user.Password.Match(input.Password)
	if err != nil {
		app.writeError(w, r, err)
		return
	}

	if !match {
		app.recordAuditLog(r, &data.AuditLog{
			UserID:     nil,
			Action:     "login",
			Resource:   "user",
			ResourceID: nil,
			Result:     "failed",
		})

		app.writeError(w, r, apperr.New(apperr.CodeInvalidCredentials))
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.writeError(w, r, err)
		return
	}

	app.recordAuditLog(r, &data.AuditLog{
		UserID:     &user.ID,
		Action:     "login",
		Resource:   "user",
		ResourceID: &user.ID,
		Result:     "success",
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication": token}, nil)
	if err != nil {
		app.writeError(w, r, err)
	}
}
