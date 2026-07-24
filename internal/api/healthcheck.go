package api

import "net/http"

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.cfg.Env,
			"version":     version,
		},
	}

	if err := app.writeJSON(w, http.StatusOK, data, nil); err != nil {
		app.logger.Error(err.Error())
	}
}
