package middleware

import "net/http"

// ResponseWriter is a wrapper for http.ResponseWriter that captures the status code
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
