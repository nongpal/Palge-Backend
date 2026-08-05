package api

import (
	"errors"
	"fmt"
	"net/http"

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
		app.logger.Error(err.Error())
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		err = app.mailer.Send(users.Email, "user_welcome.tmpl.html", users)
		if err != nil {
			app.logger.Error(err.Error())
		}
	}()

	if err := app.writeJSON(w, http.StatusAccepted, envelope{
		"user": users,
	}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
