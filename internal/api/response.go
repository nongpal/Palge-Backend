package api

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	status int
}


