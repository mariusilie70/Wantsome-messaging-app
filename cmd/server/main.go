package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"wantsome.ro/messagingapp/internal/server"
)

func main() {
	// Create a new server instance
	s := server.NewServer()

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("Shutting down the server...")
		s.Close()
		os.Exit(0)
	}()

	// Start the server
	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, s.Router()); err != nil {
		log.Fatal(err)
	}
}
