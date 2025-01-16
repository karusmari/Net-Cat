package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client represents a chat client
type Client struct {
	conn     net.Conn
	name     string
	outgoing chan string
}

func (s *Server) ClientName(conn net.Conn) (string, error) {
	scanner := bufio.NewScanner(conn)
	conn.Write([]byte(welcomeArt))
	var name string

	for scanner.Scan() {
		name = strings.TrimSpace(scanner.Text())
		if name == "" {
			conn.Write([]byte("Name cannot be empty.\n[ENTER YOUR NAME]: "))
			continue
		}
		if !s.IsValidName(name) {
			conn.Write([]byte("Name cannot contain unprintable characters.\n[ENTER YOUR NAME]: "))
			continue
		}
		if s.isNameTaken(name) {
			conn.Write([]byte("This name is already taken.\n[ENTER YOUR NAME]: "))
			continue
		}
		return name, nil
	}
	return "", fmt.Errorf("client disconnected")
}

// checks messages what client sends into server
func (s *Server) HandleIncomingMessages(client *Client) {
	//reads incoming messages from the client
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		msg := strings.TrimSpace(scanner.Text())
		msg = strings.ReplaceAll(msg, "\033", "")

		if msg == "/exit" || msg == "/quit" {
			client.conn.Write([]byte("You have left the chat. Goodbye!\n"))
			return
		}

		if msg != "" {
			//if the message is normal, we format it, broadcast it and log it
			fullMsg := FormatMessage(client.name, msg)
			logToFile(LogEntry{
				Timestamp: time.Now(),
				Username:  client.name,
				Message:   msg,
			})
			s.broadcast(fullMsg, nil)
		}
	}
}

// sends messages to certain clients after the broadcast function
func (s *Server) sendMessages(client *Client) {
	for msg := range client.outgoing {
		_, err := client.conn.Write([]byte(msg))
		if err != nil {
			break
		}
	}
}
