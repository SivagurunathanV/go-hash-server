package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHashContent_ExistingJob(t *testing.T) {
	// Seed a completed job with hash result
	StoreJob(Job{
		JobId:       "hash-test-1",
		Status:      "COMPLETED",
		HashContent: "hello world",
		HashResult:  "5eb63bbbe01eeed093cb22bb8f5acdc3",
	})

	req := httptest.NewRequest(http.MethodGet, "/hash-content/hash-test-1", nil)
	req.SetPathValue("id", "hash-test-1")
	rec := httptest.NewRecorder()

	GetHashContent(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var hashResult string
	err := json.NewDecoder(rec.Body).Decode(&hashResult)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if hashResult != "5eb63bbbe01eeed093cb22bb8f5acdc3" {
		t.Errorf("expected hash result '5eb63bbbe01eeed093cb22bb8f5acdc3', got %s", hashResult)
	}
}

func TestGetHashContent_NonExistentJob(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hash-content/no-such-job", nil)
	req.SetPathValue("id", "no-such-job")
	rec := httptest.NewRecorder()

	GetHashContent(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Should return error JSON since job doesn't exist
	var appErr AppError
	err := json.NewDecoder(rec.Body).Decode(&appErr)
	if err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if appErr.Code != 404 {
		t.Errorf("expected error code 404, got %d", appErr.Code)
	}
}
