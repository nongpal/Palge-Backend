package api

import (
	"errors"
	"fmt"
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
		Owner:   input.Owner,
		Balance: input.InitialBalance,
	}

	v := validator.New()

	data.ValidateAccount(v, account)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err = app.models.Accounts.Insert(account); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/accounts/%d", account.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"account": account}, headers)

	if err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *Application) listAccountHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := app.models.Accounts.GetAll()
	// WARN: not yet handled err returned from GetAll

	err = app.writeJSON(w, http.StatusOK, envelope{
		"account": accounts,
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

	account, err := app.models.Accounts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrAccountNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
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

	account, err := app.models.Accounts.Deposit(id, input.Amount)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrAccountNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

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

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Amount <= 0 {
		app.badRequestResponse(w, r, data.ErrInvalidAmount)
		return
	}

	account, err := app.models.Accounts.Withdraw(id, input.Amount)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrAccountNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrInsufficientBalance):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

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

	if input.Amount <= 0 {
		app.badRequestResponse(w, r, data.ErrInvalidAmount)
		return
	}

	if input.From == input.To {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, data.ErrSameAccountTransfer.Error())
		return
	}

	sender, err := app.models.Accounts.Get(input.From)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrAccountNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if sender.Balance < input.Amount {
		app.badRequestResponse(w, r, data.ErrInsufficientBalance)
		return
	}

	sender, receiver, err := app.models.Accounts.Transfer(input.From, input.To, input.Amount)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrAccountNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

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
