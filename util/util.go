package util

import (
	"fmt"
	"time"
)

// Message contains message data.
type Message struct {
	User      string
	Contents  string
	Timestamp time.Time
}

// Command contains command data.
type Command struct {
	ID         string
	OriginUser string
	Message    Message
}

// ParseTime retrieves a string from the passed time.
func ParseTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}
