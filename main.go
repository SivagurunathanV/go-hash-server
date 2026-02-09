package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"example.com/hash_server/server/handler"
)

func main() {
	handler.StartWorkerPool(runtime.NumCPU(), 100)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	httpServer := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		log.Println("Server starting on :8080")
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Shutdown failed: %v", err)
	}
	handler.StopWorkerPool()
	log.Println("Server stopped")
}
