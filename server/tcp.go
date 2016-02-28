package server

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/TransactPRO/workshop-golang-19032016/util"
)

// Listener contains TCP listener's data.
type Listener struct {
	tcpListener       *net.TCPListener
	done              chan bool
	doneAck           chan bool
	activeConnections map[string]*net.TCPConn
	cn                ConnnectionNotifier
	userName          string
}

// ConnnectionNotifier is being called on every new connection.
type ConnnectionNotifier func(string, bool)

// NewListener returns new TCP listener.
func NewListener(port int, cn ConnnectionNotifier, userName string) (l *Listener, err error) {
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
		cn:                cn,
		userName:          userName,
	}

	return
}

func (l *Listener) getListOfActiveUsers() (data []byte, err error) {
	users := []string{l.userName}
	for user := range l.activeConnections {
		users = append(users, user)
	}
	data, err = json.Marshal(users)
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

			l.cn(userName, true)

			l.activeConnections[userName] = tcpConn

			userNames, userNamesErr := l.getListOfActiveUsers()
			if userNamesErr != nil {
				log.Println("failed to generate a list of active users:" + userNamesErr.Error())
				continue
			}
			tcpConn.Write(append(userNames, '\n'))
		}
	}()
}

// SendToClients sends the message to connected clients.
func (l *Listener) SendToClients(cmd util.Command) {
	byteData, err := json.Marshal(cmd)
	if err != nil {
		log.Fatal(err)
	}
	for user, conn := range l.activeConnections {
		if user != cmd.OriginUser {
			_, connErr := conn.Write(append(byteData, '\n'))
			if connErr != nil {
				delete(l.activeConnections, user)
			}
		}
	}
}

// Stop stops active TCP listener.
func (l *Listener) Stop() {
	l.tcpListener.Close()
	l.done <- true
	<-l.doneAck
}
