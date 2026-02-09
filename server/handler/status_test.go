package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatus_ExistingJob(t *testing.T) {
	StoreJob(Job{JobId: "status-test-1", Status: "PENDING", FilePath: "/tmp/hash-server/uploads/data.txt"})

	req := httptest.NewRequest(http.MethodGet, "/status/status-test-1", nil)
	req.SetPathValue("id", "status-test-1")
	rec := httptest.NewRecorder()

	GetStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var status string
	err := json.NewDecoder(rec.Body).Decode(&status)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if status != "PENDING" {
		t.Errorf("expected PENDING, got %s", status)
	}
}

func TestGetStatus_NonExistentJob(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/status/status-missing", nil)
	req.SetPathValue("id", "status-missing")
	rec := httptest.NewRecorder()

	GetStatus(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
