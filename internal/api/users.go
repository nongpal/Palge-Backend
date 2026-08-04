package api

import (
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

	users := data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	users.Password.Set(input.Password)

	data.ValidateUser(v, &users)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Users.Insert(&users); err != nil {
		// TODO: handle duplicate email so don't fallback to server!
		app.logger.Error(err.Error())
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusAccepted, envelope{
		"name":      users.Name,
		"Activated": users.Activated,
	}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
