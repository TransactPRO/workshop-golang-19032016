package server

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/TransactPRO/workshop-golang-19032016/util"
)

// Server contains server's settings.
type Server struct {
	endpont  string
	port     int
	listener net.Listener
	messages chan util.Message
}

// New returns new server.
func New(endpont string, port int, messages chan util.Message) *Server {
	return &Server{
		endpont:  endpont,
		port:     port,
		messages: messages,
	}
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	// Parsing the reuest data.
	userName := r.FormValue("user")

	message := util.DecodeMessage(r.FormValue("msg"))

	returnError := func(msg string, statusCode int) {
		w.WriteHeader(statusCode)
		w.Write([]byte(msg))
	}

	// Returning an error if provided message data is invalid.
	if userName == "" || message == "" {
		returnError("invalid username or message string", http.StatusBadRequest)
		return
	}

	// Pushing the message in the routine because we don't want to make the client wait.
	go func(userName, message string) {
		s.messages <- util.Message{
			User:      userName,
			Contents:  message,
			Timestamp: time.Now(),
		}
	}(userName, message)

	// Telling the client that the request has been processed successfully.
	returnError("OK", http.StatusOK)
}

// Start starts the HTTP server.
func (s *Server) Start() (err error) {
	http.HandleFunc(s.endpont, s.handler)

	s.listener, err = net.Listen("tcp", ":"+strconv.Itoa(s.port))
	if err != nil {
		return
	}

	go http.Serve(s.listener, nil)

	return
}

// Stop stops the server.
func (s *Server) Stop() {
	s.listener.Close()
}
