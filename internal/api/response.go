package api

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int)  {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error ) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}

	return rw.ResponseWriter.Write(b)
}