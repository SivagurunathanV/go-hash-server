package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHashContent_ExistingJob(t *testing.T) {
	StoreJob(Job{
		JobId:      "hash-test-1",
		Status:     "COMPLETED",
		FilePath:   "/tmp/hash-server/uploads/test.txt",
		HashResult: "5eb63bbbe01eeed093cb22bb8f5acdc3",
	})

	req := httptest.NewRequest(http.MethodGet, "/hash-content/hash-test-1", nil)
	req.SetPathValue("id", "hash-test-1")
	rec := httptest.NewRecorder()

	GetHashContent(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
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
	req := httptest.NewRequest(http.MethodGet, "/hash-content/hash-missing", nil)
	req.SetPathValue("id", "hash-missing")
	rec := httptest.NewRecorder()

	GetHashContent(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
