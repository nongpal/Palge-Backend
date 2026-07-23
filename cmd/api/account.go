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

func (app *application) listAccountHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{
		"account": app.accounts,
	}, nil)

	if err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *application) showAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		// BUG: is this good enough for error handling?
		app.logger.Error(err.Error())
		return
	}

	account, err := app.GetAccountByID(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"account": account}, nil)
	if err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *application) depositHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount int64 `json:"amount"`
	}

	id, err := app.readIDParam(r)
	if err != nil {
		// BUG: is this good enough for error handling?
		app.logger.Error(err.Error())
		return
	}

	account, err := app.GetAccountByID(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.readJSON(w, r, &input)
	// TODO: not yet handled!
	account.Balance += input.Amount

	err = app.writeJSON(w, http.StatusOK, envelope{"account": account}, nil)
	if err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *application) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount int64 `json:"amount"`
	}

	id, err := app.readIDParam(r)
	if err != nil {
		// BUG: is this good enough for error handling?
		app.logger.Error(err.Error())
		return
	}

	account, err := app.GetAccountByID(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.readJSON(w, r, &input)
	// TODO: not yet handled!
	account.Balance -= input.Amount

	err = app.writeJSON(w, http.StatusOK, envelope{"account": account}, nil)
	if err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *application) transferHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		From   int64 `json:"from"`
		To     int64 `json:"to"`
		Amount int64 `json:"amount"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	sender, err := app.GetAccountByID(input.From)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	receiver, err := app.GetAccountByID(input.To)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	sender.Balance -= input.Amount
	receiver.Balance += input.Amount

	err = app.writeJSON(w, http.StatusOK, envelope{
		"transfer": map[string]*data.Account{
			"from": sender,
			"to":   receiver,
		},
	}, nil)

	if err != nil {
		app.logger.Error(err.Error())
	}
}
