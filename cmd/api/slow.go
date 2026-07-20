package main

import (
	"net/http"
	"time"
)

func (app *application) slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(15 * time.Second)

	data := envelope{
		"message": "finished",
	}

	if err := app.writeJSON(w, http.StatusOK, data, nil); err != nil {
		app.logger.Error(err.Error())
	}
}
