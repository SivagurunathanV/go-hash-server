package handler

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /status/{id}", GetStatus)
	mux.HandleFunc("POST /upload", ProcessFile)
	mux.HandleFunc("GET /hash-content/{id}", GetHashContent)
	mux.HandleFunc("POST /upload_log", ProcessLogFile)
}
