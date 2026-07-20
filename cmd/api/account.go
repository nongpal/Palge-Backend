package main

import (
	"net/http"

	"github.com/nongpal/Palge-Backend/internal/data"
)

func (app *application) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Owner          string `json:"owner"`
		InitialBalance int64  `json:"initial_balance"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	account := &data.Account{
		ID:      app.nextID,
		Owner:   input.Owner,
		Balance: input.InitialBalance,
	}

	app.accounts = append(app.accounts, account)
	app.nextID++

	err = app.writeJSON(w, http.StatusCreated, envelope{
		"account": account,
	}, nil)

	if err != nil {
		app.logger.Error(err.Error())
	}
}
