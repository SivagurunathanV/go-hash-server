package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

func writeErrorResponse(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		json.NewEncoder(w).Encode(appErr)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}
}

type Job struct {
	JobId      string
	Status     string
	FilePath   string
	HashResult string
}

type AppError struct {
	Code    int
	JobId   string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

var jobStore map[string]Job = make(map[string]Job)
var lock sync.RWMutex

func StartJob(job Job) (Job, error) {
	job, err := StoreJob(job)
	if err != nil {
		return Job{}, err
	}
	EnqueueJob(job.JobId)
	return job, nil
}

func StoreJob(job Job) (Job, error) {
	lock.Lock()
	jobStore[job.JobId] = job
	lock.Unlock()
	return job, nil
}

// readJobLocked is the internal lookup — caller must hold the lock.
func readJobLocked(jobId string) (Job, error) {
	val, exists := jobStore[jobId]
	if !exists {
		return Job{}, &AppError{JobId: jobId, Code: 404, Message: fmt.Sprintf("Job %s Not found", jobId)}
	}
	return val, nil
}

func ReadJob(jobId string) (Job, error) {
	lock.RLock()
	defer lock.RUnlock()
	return readJobLocked(jobId)
}

func UpdateJob(job Job) (Job, error) {
	lock.Lock()
	defer lock.Unlock()
	jobStore[job.JobId] = job
	return job, nil
}

func updateJobResult(jobId string, status string, hashResult string) error {
	lock.Lock()
	defer lock.Unlock()
	val, err := readJobLocked(jobId)
	if err != nil {
		return err
	}
	val.Status = status
	val.HashResult = hashResult
	jobStore[jobId] = val
	return nil
}
