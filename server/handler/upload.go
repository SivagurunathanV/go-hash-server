package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const allowedBaseDir = "/tmp/hash-server/uploads"

type UploadRequest struct {
	FilePath string `json:"filepath"`
}

func newJobId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func validateFilePath(path string) error {
	cleaned := filepath.Clean(path)
	if !strings.HasPrefix(cleaned, allowedBaseDir) {
		return &AppError{Code: 400, Message: "file path must be within " + allowedBaseDir}
	}
	if _, err := os.Stat(cleaned); err != nil {
		return &AppError{Code: 400, Message: "file does not exist: " + cleaned}
	}
	return nil
}

func ProcessFile(w http.ResponseWriter, r *http.Request) {
	var input UploadRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := validateFilePath(input.FilePath); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeErrorResponse(w, err)
		return
	}
	fmt.Printf("Content: %s", input.FilePath)
	newJob := Job{
		JobId:    newJobId(),
		Status:   "PENDING",
		FilePath: filepath.Clean(input.FilePath),
	}
	storedJob, err := StartJob(newJob)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(storedJob)
}
