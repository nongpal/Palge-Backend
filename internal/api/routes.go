package api

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("GET /v1/slow", app.slowHandler)
	mux.HandleFunc("POST /v1/accounts", app.createAccountHandler)
	mux.HandleFunc("GET /v1/accounts", app.listAccountHandler)
	mux.HandleFunc("GET /v1/accounts/{id}", app.showAccountHandler)
	mux.HandleFunc("POST /v1/accounts/{id}/deposit", app.depositHandler)
	mux.HandleFunc("POST /v1/accounts/{id}/withdraw", app.withdrawHandler)
	mux.HandleFunc("POST /v1/transfers", app.transferHandler)
	mux.HandleFunc("POST /v1/users", app.registerUserHandler)
	mux.HandleFunc("PUT /v1/users/activated", app.activateUserHandler)

	return mux
}
