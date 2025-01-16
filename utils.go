package main

import (
	"fmt"
	"time"
)


func TimeStamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// FormatMessage is for chat messages
func FormatMessage(username, message string) string {
	return fmt.Sprintf("[%s][%s]:%s\n", TimeStamp(), username, message)
}

// FormatSystemMessage is for system messages
func FormatSystemMessage(username, action string) string {
	return fmt.Sprintf("[%s] %s %s\n", TimeStamp(), username, action)
}

func (s *Server) IsValidName(name string) bool {
	for _, char := range name {
		if char < 32 || char > 126 {
			return false
		}
	}
	return true
}