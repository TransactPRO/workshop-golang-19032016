package server

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
)

// Listener contains TCP listener's data.
type Listener struct {
	tcpListener       *net.TCPListener
	done              chan bool
	doneAck           chan bool
	activeConnections map[string]*net.TCPConn
}

// NewListener returns new TCP listener.
func NewListener(port int) (l *Listener, err error) {
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		return
	}

	var tcpListener *net.TCPListener
	tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}

	l = &Listener{
		tcpListener:       tcpListener,
		done:              make(chan bool),
		doneAck:           make(chan bool),
		activeConnections: make(map[string]*net.TCPConn),
	}

	return
}

// Start makes the TCP listener to start accepting incoming connections.
func (l *Listener) Start() {
	go func() {
		var closed bool
		for {
			// Waiting for a new client to connect.
			tcpConn, tcpErr := l.tcpListener.AcceptTCP()

			select {
			case <-l.done:
				closed = true
			default:
			}
			if closed {
				l.doneAck <- true
				break
			}

			if tcpErr != nil {
				log.Println("failed to establish TCP client connection")
				continue
			}

			result, resErr := bufio.NewReader(tcpConn).ReadString('\n')
			if resErr != nil {
				log.Println("failed to process data from master:" + resErr.Error())
				continue
			}

			userName := strings.Replace(result, "\n", "", -1)
			log.Println(userName, "connected")

			l.activeConnections[userName] = tcpConn
		}
	}()
}

// Stop stops active TCP listener.
func (l *Listener) Stop() {
	l.tcpListener.Close()
	l.done <- true
	<-l.doneAck
}
