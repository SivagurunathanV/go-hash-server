package handler

import (
	"fmt"
	"sync"
)

type Job struct {
	JobId       string
	Status      string
	HashContent string
	HashResult  string
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

func StoreJob(job Job) (Job, error) {
	lock.Lock()
	jobStore[job.JobId] = job
	lock.Unlock()
	return job, nil
}

func ReadJob(jobId string) (Job, error) {
	lock.RLock()
	defer lock.RUnlock()
	val, exists := jobStore[jobId]
	if !exists {
		return Job{}, &AppError{JobId: jobId, Code: 404, Message: fmt.Sprintf("Job %s Not found", jobId)}
	}
	return val, nil
}

func updateJob(jobId string, status string) (Job, error) {
	val, err := ReadJob(jobId)
	if err != nil {
		return Job{}, err
	}
	val.Status = status
	job, err := StoreJob(val)
	if err != nil {
		return Job{}, err
	}
	return job, nil
}
