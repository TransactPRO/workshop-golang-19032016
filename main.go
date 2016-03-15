package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TransactPRO/workshop-golang-19032016/client"
	"github.com/TransactPRO/workshop-golang-19032016/server"
	"github.com/TransactPRO/workshop-golang-19032016/ui"
	"github.com/TransactPRO/workshop-golang-19032016/util"
)

var (
	userName   = flag.String("u", "", "master host")
	masterMode = flag.Bool("m", false, "run in master mode")
	httpPort   = flag.Int("http-port", 8081, "master's HTTP port")
	masterHost = flag.String("master-host", "127.0.0.1", "master host")
	tcpPort    = flag.Int("tcp-port", 8082, "master's TCP port")
)

var (
	gui             *ui.UI
	textBoxMessages = make(chan string)
	srv             *server.Server
	clientsMessages = make(chan util.Message)
	l               *server.Listener
)

// shutdown stops all the running services and terminates the process.
func shutdown() {
	log.Println("Shutting down..")

	if srv != nil {
		srv.Stop()
	}

	if l != nil {
		l.Stop()
	}

	os.Exit(0)
}

func processMyMessages() {
	for msgStr := range textBoxMessages {
		msg := util.Message{
			User:      *userName,
			Contents:  msgStr,
			Timestamp: time.Now(),
		}
		gui.WriteToView(ui.ChatView, fmt.Sprintf("[%s] %s: %s", util.ParseTime(msg.Timestamp), msg.User, msg.Contents))

		if !*masterMode {
			client.SendMessageToMaster(msg, *masterHost, *httpPort)
		} else {
			//send to clients
		}
	}
}

func processClientsMessages() {
	for msg := range clientsMessages {
		gui.WriteToView(ui.ChatView, fmt.Sprintf("[%s] %s: %s", util.ParseTime(msg.Timestamp), msg.User, msg.Contents))
	}
}

func processNewUser(newUser string, notifyClients bool) {
	gui.WriteToView(ui.UsersView, newUser)
	if notifyClients {
		l.SendToClients(util.Command{
			ID:         "USER",
			OriginUser: newUser,
		})
	}
}

func main() {
	flag.Parse()

	if *userName == "" {
		log.Fatal("provide the username!")
	}

	var err error

	gui, err = ui.DeployGUI(shutdown, textBoxMessages)
	if err != nil {
		log.Fatal(err)
	}

	if *masterMode {
		// Create a new server with the desired parameters.
		srv = server.New("/", *httpPort, clientsMessages)
		// Start the server (initialize the TCP listener).
		err = srv.Start()
		if err != nil {
			log.Fatal(err)
		}

		// Create a new TCP listener.
		l, err = server.NewListener(*tcpPort, *userName, processNewUser)
		if err != nil {
			log.Fatal(err)
		}
		// Start the listener.
		l.Start()
	} else {

	}

	go processMyMessages()
	go processClientsMessages()

	select {}
}
