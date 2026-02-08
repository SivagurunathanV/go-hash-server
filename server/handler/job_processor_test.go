package handler

import (
	"errors"
	"testing"
)

func TestStoreJob(t *testing.T) {
	newJob := Job{
		JobId:       "store-test-1",
		Status:      "PENDING",
		HashContent: "test content",
	}

	job, err := StoreJob(newJob)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.JobId != "store-test-1" {
		t.Errorf("expected JobId 'store-test-1', got %s", job.JobId)
	}
	if job.Status != "PENDING" {
		t.Errorf("expected status PENDING, got %s", job.Status)
	}
	if job.HashContent != "test content" {
		t.Errorf("expected HashContent 'test content', got %s", job.HashContent)
	}
}

func TestReadJob_Existing(t *testing.T) {
	StoreJob(Job{JobId: "read-test-1", Status: "PENDING", HashContent: "data"})

	job, err := ReadJob("read-test-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.JobId != "read-test-1" {
		t.Errorf("expected JobId 'read-test-1', got %s", job.JobId)
	}
}

func TestReadJob_NonExistent(t *testing.T) {
	_, err := ReadJob("does-not-exist")
	if err == nil {
		t.Fatal("expected error for non-existent job, got nil")
	}

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *AppError, got %T", err)
	}
	if appErr.Code != 404 {
		t.Errorf("expected error code 404, got %d", appErr.Code)
	}
	if appErr.JobId != "does-not-exist" {
		t.Errorf("expected JobId 'does-not-exist', got %s", appErr.JobId)
	}
}

func TestUpdateJob(t *testing.T) {
	StoreJob(Job{JobId: "update-test-1", Status: "PENDING", HashContent: "data"})

	job, err := updateJob("update-test-1", "COMPLETED")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got %s", job.Status)
	}

	// Verify persisted
	stored, _ := ReadJob("update-test-1")
	if stored.Status != "COMPLETED" {
		t.Errorf("expected stored status COMPLETED, got %s", stored.Status)
	}
}

func TestUpdateJob_NonExistent(t *testing.T) {
	_, err := updateJob("no-such-job", "COMPLETED")
	if err == nil {
		t.Fatal("expected error for non-existent job, got nil")
	}

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *AppError, got %T", err)
	}
	if appErr.Code != 404 {
		t.Errorf("expected error code 404, got %d", appErr.Code)
	}
}
