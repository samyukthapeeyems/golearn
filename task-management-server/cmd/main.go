package main

import (
	"log"
	"net/http"
	"os"

	"example.com/task-management-server/internal/route"
)

func main() {
	serverPort := os.Getenv("SERVERPORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	router := route.SetupRoutes()

	log.Printf("server starting on port %s\n", serverPort)
	err := http.ListenAndServe(":"+serverPort, router)
	if err != nil {
		log.Fatal(err)
	}
}
