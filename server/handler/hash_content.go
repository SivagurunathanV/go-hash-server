package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func GetHashContent(w http.ResponseWriter, r *http.Request) {
	jobId := r.PathValue("id")
	fmt.Printf("Input request %s", jobId)
	w.Header().Set("Content-Type", "application/json")
	job, err := ReadJob(jobId)
	if err != nil {
		var appErr *AppError
		if errors.As(err, &appErr) {
			w.WriteHeader(appErr.Code)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		writeErrorResponse(w, err)
		return
	}
	json.NewEncoder(w).Encode(job.HashResult)
}
