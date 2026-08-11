package api

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("GET /v1/slow", app.slowHandler)

	mux.HandleFunc("POST /v1/accounts", app.requiredActivatedUser(app.createAccountHandler))
	mux.HandleFunc("GET /v1/accounts", app.requiredActivatedUser(app.listAccountHandler))
	mux.HandleFunc("GET /v1/accounts/{id}", app.requiredActivatedUser(app.showAccountHandler))
	mux.HandleFunc("POST /v1/accounts/{id}/deposit", app.requiredActivatedUser(app.depositHandler))
	mux.HandleFunc("POST /v1/accounts/{id}/withdraw", app.requiredActivatedUser(app.withdrawHandler))
	mux.HandleFunc("POST /v1/transfers", app.requiredActivatedUser(app.transferHandler))

	mux.HandleFunc("POST /v1/users", app.registerUserHandler)
	mux.HandleFunc("PUT /v1/users/activated", app.activateUserHandler)
	mux.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.authenticate(mux)
}
