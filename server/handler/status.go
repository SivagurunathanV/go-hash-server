package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetStatus(writer http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing status")
	writer.Header().Set("Content-Type", "application/json")
	jobId := r.PathValue("id")
	jobStatus, err := ReadJob(jobId)
	if err != nil {
		json.NewEncoder(writer).Encode(err)
		return
	}
	json.NewEncoder(writer).Encode(jobStatus)
}
