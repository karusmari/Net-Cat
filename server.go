package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	port       int
	maxClients int
	clients    map[*Client]bool
	messages   []string
	mutex      sync.Mutex
}

// creating a new server
func NewServer(port, maxClients int) (*Server, error) {
	return &Server{
		port:       port,
		maxClients: maxClients,
		clients:    make(map[*Client]bool),
		messages:   make([]string, 0), //creating an empty slice of strings
	}, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		s.mutex.Lock()
		if len(s.clients) >= s.maxClients {
			s.mutex.Unlock()
			conn.Write([]byte("Chat is full. Try again later.\n"))
			conn.Close()
			continue
		}
		s.mutex.Unlock()

		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	//get the client's name
	clientName, err := s.ClientName(conn)
	if err != nil {
		return
	}

	client := s.AddClient(conn, clientName) //add client to the server
	defer s.DisconnectClient(client)

	go s.sendMessages(client)
	s.HandleIncomingMessages(client)
}

func (s *Server) WelcomeClient(client *Client) {
	welcomeMsg := FormatSystemMessage(client.name, "has joined the chat")
	s.broadcast(welcomeMsg, client)
	s.sendHistory(client)
}

func (s *Server) AddClient(conn net.Conn, name string) *Client {
	//create a new client
	client := &Client{
		conn:     conn,
		name:     name,
		outgoing: make(chan string, 100),
	}

	//add client to the server
	s.mutex.Lock()
	s.clients[client] = true
	s.mutex.Unlock()

	log.Printf("%s has connected.\n", client.name)

	logToFile(LogEntry{
		Timestamp: time.Now(),
		Username:  client.name,
		Action:    "connected to chat",
	})

	go s.sendMessages(client)
	s.sendHistory(client)

	welcomeMsg := FormatSystemMessage(client.name, "has joined the chat")
	s.broadcast(welcomeMsg, client)

	return client
}

func (s *Server) DisconnectClient(client *Client) {
	s.mutex.Lock()
	delete(s.clients, client)
	close(client.outgoing)
	s.mutex.Unlock()

	disconnectMsg := FormatSystemMessage(client.name, "has left the chat.")
	s.broadcast(disconnectMsg, client)

	log.Printf("%s has disconnected.\n", client.name)

	logToFile(LogEntry{
		Timestamp: time.Now(),
		Username:  client.name,
		Action:    "disconnected from chat",
	})
}

// this function will send messages to the client
func (s *Server) broadcast(message string, except *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.messages = append(s.messages, message)

	for client := range s.clients {
		if client != except {
			select {
			case client.outgoing <- message:
			default:
				// Channel is full, skip this client and continue to the next one
			}
		}
	}
}

func (s *Server) sendHistory(client *Client) {
	s.mutex.Lock()
	history := make([]string, len(s.messages))
	copy(history, s.messages)
	s.mutex.Unlock()

	for _, msg := range history {
		client.outgoing <- msg
	}
}

func (s *Server) isNameTaken(name string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for client := range s.clients {
		if strings.EqualFold(client.name, name) { // case-insensitive comparison
			return true
		}
	}
	return false
}
