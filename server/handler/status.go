package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func GetStatus(writer http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing status")
	writer.Header().Set("Content-Type", "application/json")
	jobId := r.PathValue("id")
	job, err := ReadJob(jobId)
	if err != nil {
		var appErr *AppError
		if errors.As(err, &appErr) {
			writer.WriteHeader(appErr.Code)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		writeErrorResponse(writer, err)
		return
	}
	json.NewEncoder(writer).Encode(job.Status)
}
