package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createMultipartRequest(t *testing.T, fieldName, fileName, content string) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write([]byte(content))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload_log", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestProcessLogFile_ValidUpload(t *testing.T) {
	os.MkdirAll(allowedBaseDir, 0755)
	req := createMultipartRequest(t, "file", "test-log.txt", "hello world")
	rec := httptest.NewRecorder()

	ProcessLogFile(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var job Job
	err := json.NewDecoder(rec.Body).Decode(&job)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if job.JobId == "" {
		t.Error("expected non-empty JobId")
	}
	if job.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got %s", job.Status)
	}
	if job.HashResult == "" {
		t.Error("expected non-empty HashResult")
	}
	if job.FilePath == "" {
		t.Error("expected non-empty FilePath")
	}
	t.Cleanup(func() { os.Remove(job.FilePath) })
}

func TestProcessLogFile_JobStoredAndQueryable(t *testing.T) {
	os.MkdirAll(allowedBaseDir, 0755)
	req := createMultipartRequest(t, "file", "test-stored.txt", "store me")
	rec := httptest.NewRecorder()

	ProcessLogFile(rec, req)

	var respJob Job
	json.NewDecoder(rec.Body).Decode(&respJob)
	t.Cleanup(func() { os.Remove(respJob.FilePath) })

	stored, err := ReadJob(respJob.JobId)
	if err != nil {
		t.Fatalf("job %s not found in store: %v", respJob.JobId, err)
	}
	if stored.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got %s", stored.Status)
	}
	if stored.HashResult != respJob.HashResult {
		t.Errorf("expected HashResult %s, got %s", respJob.HashResult, stored.HashResult)
	}
}

func TestProcessLogFile_MissingFileField(t *testing.T) {
	req := createMultipartRequest(t, "wrong_field", "test.txt", "data")
	rec := httptest.NewRecorder()

	ProcessLogFile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestProcessLogFile_EmptyFileContent(t *testing.T) {
	os.MkdirAll(allowedBaseDir, 0755)
	req := createMultipartRequest(t, "file", "empty.txt", "")
	rec := httptest.NewRecorder()

	ProcessLogFile(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var job Job
	json.NewDecoder(rec.Body).Decode(&job)
	t.Cleanup(func() { os.Remove(job.FilePath) })

	if job.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got %s", job.Status)
	}
	// MD5 of empty string is d41d8cd98f00b204e9800998ecf8427e
	if job.HashResult != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("expected empty MD5 hash, got %s", job.HashResult)
	}
}
