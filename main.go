package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	defaultPort = 8989
	maxClients  = 10
)

func main() {
	port := defaultPort
	if len(os.Args) > 1 {
		if len(os.Args) > 2 {
			fmt.Println("[USAGE]: ./TCPChat $port")
			os.Exit(1)
		}

		p, err := strconv.Atoi(os.Args[1])

		if err != nil || p < 1024 || p > 65535 {
			fmt.Println("[USAGE]: ./TCPChat $port")
			os.Exit(1)
		}
		port = p
	}

	server, err := NewServer(port, maxClients)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Printf("Starting TCP Chat server on port %d...", port)
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
