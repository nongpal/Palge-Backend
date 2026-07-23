package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("GET /v1/slow", app.slowHandler)
	mux.HandleFunc("POST /v1/accounts", app.createAccountHandler)
	mux.HandleFunc("GET /v1/accounts", app.listAccountHandler)
	mux.HandleFunc("GET /v1/accounts/{id}", app.showAccountHandler)
	mux.HandleFunc("POST /v1/accounts/{id}/deposit", app.depositHandler)
	mux.HandleFunc("POST /v1/accounts/{id}/withdraw", app.withdrawHandler)

	return mux
}
