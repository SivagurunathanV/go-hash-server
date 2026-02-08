package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatus_ExistingJob(t *testing.T) {
	StoreJob(Job{JobId: "status-test-1", Status: "PENDING", HashContent: "data"})

	req := httptest.NewRequest(http.MethodGet, "/status/status-test-1", nil)
	req.SetPathValue("id", "status-test-1")
	rec := httptest.NewRecorder()

	GetStatus(rec, req)

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
	if job.Status != "PENDING" {
		t.Errorf("expected PENDING, got %s", job.Status)
	}
}

func TestGetStatus_NonExistentJob(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/status/does-not-exist", nil)
	req.SetPathValue("id", "does-not-exist")
	rec := httptest.NewRecorder()

	GetStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var appErr AppError
	err := json.NewDecoder(rec.Body).Decode(&appErr)
	if err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if appErr.Code != 404 {
		t.Errorf("expected error code 404, got %d", appErr.Code)
	}
}
