package util

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

// Message contains message data.
type Message struct {
	User      string
	Contents  string
	Timestamp time.Time
}

// ParseTime retrieves a string from the passed time.
func ParseTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

// EncodeMessage encodes the message to base64.
func EncodeMessage(msg string) string {
	return base64.StdEncoding.EncodeToString([]byte(msg))
}

// DecodeMessage decodes the message from base64.
func DecodeMessage(msg string) (decoded string) {
	sDec, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		log.Println(err)
	}
	return string(sDec)
}
