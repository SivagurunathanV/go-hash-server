package handler

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var jobChan chan string
var wg sync.WaitGroup
var pool sync.Pool

func StartWorkerPool(workerCount int, bufferSize int) {
	pool = sync.Pool{
		New: func() any {
			return make([]byte, 32*1024)
		},
	}
	jobChan = make(chan string, bufferSize)
	for i := range workerCount {
		wg.Add(1)
		go worker(i)
	}
}

func worker(id int) {
	defer wg.Done()
	for jobId := range jobChan {
		fmt.Printf("Worker %d picked up job %s\n", id, jobId)
		// TODO: open temp file, stream through hash, update job
		processJob(jobId)
	}
}

func processJob(jobId string) {
	job, err := ReadJob(jobId)
	if err != nil {
		log.Printf("Job %s not found: %v", jobId, err)
		return
	}
	result, err := computeHash(job.FilePath)
	if err != nil {
		log.Printf("Job %s failed during hashing: %v", jobId, err)
		updateJobResult(jobId, "FAILED", "")
		return
	}
	updateJobResult(jobId, "COMPLETED", result)
}

func computeHash(path string) (string, error) {
	// read file
	hasher := md5.New()
	file, err := os.Open(path)
	if err != nil {
		log.Println("File error", err)
		return "", err
	}
	defer file.Close()
	buf := pool.Get().([]byte)
	defer pool.Put(buf)
	count, err := io.CopyBuffer(hasher, file, buf)
	if err != nil {
		log.Println("Failed while copying", err)
		return "", err
	}
	log.Printf("Copied count %d", count)
	content := fmt.Sprintf("%x", hasher.Sum(nil))
	return content, nil
}

func EnqueueJob(jobId string) {
	jobChan <- jobId
}

func StopWorkerPool() {
	close(jobChan)
	wg.Wait()
}
