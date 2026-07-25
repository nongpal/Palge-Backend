package api

import (
	"net/http"

	"github.com/nongpal/Palge-Backend/internal/data"
	"github.com/nongpal/Palge-Backend/internal/validator"
)

func (app *Application) createAccountHandler(w http.ResponseWriter, r *http.Request) {
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

	v := validator.New()

	data.ValidateAccount(v, account)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
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

func (app *Application) listAccountHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{
		"account": app.accounts,
	}, nil)

	if err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *Application) showAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
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

func (app *Application) depositHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount int64 `json:"amount"`
	}

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	account, err := app.GetAccountByID(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateAmount(v, input.Amount)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	account.Balance += input.Amount

	err = app.writeJSON(w, http.StatusOK, envelope{"account": account}, nil)
	if err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *Application) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount int64 `json:"amount"`
	}

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	account, err := app.GetAccountByID(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Amount <= 0 {
		app.badRequestResponse(w, r, data.ErrInvalidAmount)
		return
	}

	if account.Balance < input.Amount {
		app.badRequestResponse(w, r, data.ErrInsufficientBalance)
		return
	}

	account.Balance -= input.Amount

	err = app.writeJSON(w, http.StatusOK, envelope{"account": account}, nil)
	if err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *Application) transferHandler(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	receiver, err := app.GetAccountByID(input.To)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if input.Amount <= 0 {
		app.badRequestResponse(w, r, data.ErrInvalidAmount)
		return
	}

	if sender.Balance < input.Amount {
		app.badRequestResponse(w, r, data.ErrInsufficientBalance)
		return
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
