package main

import (
	"fmt"
	"os"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Username  string
	Action    string
	Message   string
}

// this function logs the chat history to a file
func logToFile(entry LogEntry) error {
	file, err := os.OpenFile("chat_history.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var logLine string
	if entry.Message != "" {
		logLine = fmt.Sprintf("[%s][%s]: %s\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Username,
			entry.Message)
	} else {
		logLine = fmt.Sprintf("[%s] %s: %s\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Username,
			entry.Action)
	}

	// write the logLine to the file
	if _, err := file.WriteString(logLine); err != nil {
		return err
	}

	return nil
}
