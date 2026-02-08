package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProcessFile_ValidInput(t *testing.T) {
	body := `{"content":"hello world"}`
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(body))
	rec := httptest.NewRecorder()

	ProcessFile(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var job Job
	err := json.NewDecoder(rec.Body).Decode(&job)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if job.JobId == "" {
		t.Error("expected non-empty JobId in response")
	}
	if job.Status != "PENDING" {
		t.Errorf("expected status PENDING, got %s", job.Status)
	}
	if job.HashContent != "hello world" {
		t.Errorf("expected HashContent 'hello world', got %s", job.HashContent)
	}
}

func TestProcessFile_StoresJob(t *testing.T) {
	body := `{"content":"hash me"}`
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(body))
	rec := httptest.NewRecorder()

	ProcessFile(rec, req)

	var respJob Job
	json.NewDecoder(rec.Body).Decode(&respJob)

	stored, err := ReadJob(respJob.JobId)
	if err != nil {
		t.Fatalf("job %s not found in store: %v", respJob.JobId, err)
	}
	if stored.Status != "PENDING" {
		t.Errorf("expected status PENDING, got %s", stored.Status)
	}
	if stored.HashContent != "hash me" {
		t.Errorf("expected HashContent 'hash me', got %s", stored.HashContent)
	}
}

func TestProcessFile_InvalidJSON(t *testing.T) {
	body := `not valid json`
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(body))
	rec := httptest.NewRecorder()

	ProcessFile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestProcessFile_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	rec := httptest.NewRecorder()

	ProcessFile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
