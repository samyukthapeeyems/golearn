package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Started server on port 8888")

	// Create a context for graceful shutdown
	shutdownCtx, shutdown := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Listen for OS interrupt or terminate signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Shutdown signal received")
		shutdown()
	}()

	go func() {
		// Wait for the shutdown signal and stop accepting new connections
		<-shutdownCtx.Done()
		listener.Close() // Close the listener to stop accepting new connections
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Check if the error is due to the listener being closed
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				break
			}
			log.Printf("unable to accept connection: %s", err.Error())
			continue
		}

		// Handle new client connection
		wg.Add(1)
		go func(conn net.Conn) {
			defer wg.Done()
			s.newClient(conn)
		}(conn)
	}

	// Wait for all connections to be handled
	wg.Wait()
	log.Println("Server has shut down gracefully")
}
