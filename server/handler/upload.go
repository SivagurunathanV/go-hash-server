package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type UploadRequest struct {
	Content string `json:"content"`
}

func ProcessFile(w http.ResponseWriter, r *http.Request) {
	var input UploadRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	fmt.Printf("Content: %s", input.Content)
	jobId := uuid.NewString()
	newJob := Job{
		JobId:       jobId,
		Status:      "PENDING",
		HashContent: input.Content,
	}
	storedJob, err := StoreJob(newJob)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(storedJob)
}
