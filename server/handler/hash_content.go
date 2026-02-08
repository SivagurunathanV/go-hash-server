package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetHashContentRequest struct {
	JobId string
}

func GetHashContent(w http.ResponseWriter, r *http.Request) {
	input := GetHashContentRequest{
		JobId: r.PathValue("id"),
	}
	fmt.Printf("Input request %s", input.JobId)
	w.Header().Set("Content-Type", "application/json")
	job, err := ReadJob(input.JobId)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(job.HashResult)
}
