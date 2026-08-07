package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/nongpal/Palge-Backend/internal/data"
	"github.com/nongpal/Palge-Backend/internal/validator"
)

func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	users := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	if err := users.Password.Set(input.Password); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data.ValidateUser(v, users)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Users.Insert(users); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(users.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          users.ID,
		}
		err = app.mailer.Send(users.Email, "user_welcome.tmpl.html", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	if err := app.writeJSON(w, http.StatusAccepted, envelope{
		"user": users,
	}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
