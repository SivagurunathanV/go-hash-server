package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestFile(t *testing.T, content string) string {
	t.Helper()
	os.MkdirAll(allowedBaseDir, 0755)
	f, err := os.CreateTemp(allowedBaseDir, "test-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestMain(m *testing.M) {
	StartWorkerPool(2, 10)
	code := m.Run()
	StopWorkerPool()
	os.Exit(code)
}

func TestProcessFile_ValidInput(t *testing.T) {
	path := setupTestFile(t, "hello world")
	body := `{"filepath":"` + path + `"}`
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
	if job.FilePath != path {
		t.Errorf("expected FilePath '%s', got %s", path, job.FilePath)
	}
}

func TestProcessFile_StoresJob(t *testing.T) {
	path := setupTestFile(t, "hash me")
	body := `{"filepath":"` + path + `"}`
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
	if stored.FilePath != path {
		t.Errorf("expected FilePath '%s', got %s", path, stored.FilePath)
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

func TestProcessFile_PathOutsideAllowed(t *testing.T) {
	body := `{"filepath":"/etc/passwd"}`
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(body))
	rec := httptest.NewRecorder()

	ProcessFile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestProcessFile_PathTraversal(t *testing.T) {
	body := `{"filepath":"` + filepath.Join(allowedBaseDir, "../../etc/passwd") + `"}`
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(body))
	rec := httptest.NewRecorder()

	ProcessFile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestProcessFile_FileDoesNotExist(t *testing.T) {
	body := `{"filepath":"` + filepath.Join(allowedBaseDir, "nonexistent.txt") + `"}`
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(body))
	rec := httptest.NewRecorder()

	ProcessFile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
