package gohashserver

import (
	"log"
	"net/http"

	"example.com/hash_server/server/handler"
)

func main() {
	server := http.NewServeMux()
	handler.RegisterRoutes(server)
	log.Fatal(http.ListenAndServe(":8080", server))
}
