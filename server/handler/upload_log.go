package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func ProcessLogFile(w http.ResponseWriter, r *http.Request) {
	src, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer src.Close()

	os.MkdirAll(allowedBaseDir, 0755)
	destPath := filepath.Join(allowedBaseDir, header.Filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "failed to create destination file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	hasher := md5.New()
	teeReader := io.TeeReader(src, destFile)

	buf := pool.Get().([]byte)
	defer pool.Put(buf)

	_, err = io.CopyBuffer(hasher, teeReader, buf)
	if err != nil {
		http.Error(w, "failed to process file", http.StatusInternalServerError)
		return
	}

	hashResult := fmt.Sprintf("%x", hasher.Sum(nil))

	job := Job{
		JobId:      newJobId(),
		Status:     "COMPLETED",
		FilePath:   destPath,
		HashResult: hashResult,
	}
	StoreJob(job)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}
