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

	if err = app.models.Accounts.Insert(r.Context(), account); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/accounts/%d", account.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"account": account}, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) listAccountHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := app.models.Accounts.GetAll(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{
		"account": accounts,
	}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) showAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	account, err := app.models.Accounts.Get(r.Context(), id)
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
		app.serverErrorResponse(w, r, err)
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

	account, err := app.models.Accounts.Deposit(r.Context(), id, input.Amount)
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
		app.serverErrorResponse(w, r, err)
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

	v := validator.New()
	if data.ValidateAmount(v, input.Amount); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	account, err := app.models.Accounts.Withdraw(r.Context(), id, input.Amount)
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
		app.serverErrorResponse(w, r, err)
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

	v := validator.New()
	data.ValidateAmount(v, input.Amount)
	v.Check(input.From > 0, "from", "must be a valid account ID")
	v.Check(input.To > 0, "to", "must be a valid account ID")
	v.Check(input.From != input.To, "to", "cannot transfer to the same account")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	sender, receiver, err := app.models.Accounts.Transfer(r.Context(), input.From, input.To, input.Amount)
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

	err = app.writeJSON(w, http.StatusOK, envelope{
		"transfer": map[string]*data.Account{
			"from": sender,
			"to":   receiver,
		},
	}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
