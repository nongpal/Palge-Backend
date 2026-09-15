package api

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("GET /v1/slow", app.slowHandler)

	mux.HandleFunc("POST /v1/accounts", app.requirePermission("accounts:create", app.createAccountHandler))
	mux.HandleFunc("GET /v1/accounts", app.requirePermission("accounts:read", app.listAccountHandler))
	mux.HandleFunc("GET /v1/accounts/{id}", app.requirePermission("accounts:read", app.showAccountHandler))
	mux.HandleFunc("POST /v1/accounts/{id}/deposit", app.requirePermission("accounts:deposit", app.depositHandler))
	mux.HandleFunc("POST /v1/accounts/{id}/withdraw", app.requirePermission("accounts:withdraw", app.withdrawHandler))
	mux.HandleFunc("POST /v1/transfers", app.requirePermission("accounts:transfer", app.transferHandler))

	mux.HandleFunc("POST /v1/users", app.registerUserHandler)
	mux.HandleFunc("PUT /v1/users/activated", app.activateUserHandler)

	mux.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.requestID(
		app.requestLogger(
			app.middlewareRateLimit(
				app.authenticate(mux),
			),
		),
	)
}
