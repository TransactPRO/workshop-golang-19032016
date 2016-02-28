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
	masterHost = flag.String("master-host", "127.0.0.1", "master host")
	httpPort   = flag.Int("http-port", 8081, "master's HTTP port")
	tcpPort    = flag.Int("tcp-port", 8082, "master's TCP port")
)

var (
	gui        *ui.UI
	srv        *server.Server
	l          *server.Listener
	c          *client.Client
	msgCh      = make(chan util.Message)
	myMsgCh    = make(chan string)
	newUsersCh = make(chan string)
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

	if c != nil {
		c.Close()
	}

	os.Exit(0)
}

func writeToChatView(msg util.Message) {
	gui.WriteToView(ui.ChatView, fmt.Sprintf("[%s] %s: %s", util.ParseTime(msg.Timestamp), msg.User, msg.Contents))
}

func processClientsMessages() {
	for msg := range msgCh {
		writeToChatView(msg)
		if *masterMode {
			l.SendToClients(util.Command{
				ID:         "MSG",
				OriginUser: msg.User,
				Message:    msg,
			})
		}
	}
}

func processMyMessages() {
	for msgStr := range myMsgCh {
		msg := util.Message{
			User:      *userName,
			Contents:  msgStr,
			Timestamp: time.Now(),
		}
		writeToChatView(msg)
		if *masterMode {
			l.SendToClients(util.Command{
				ID:         "MSG",
				OriginUser: *userName,
				Message:    msg,
			})
		} else {
			c.SendMessageToMaster(msg, *masterHost, *httpPort)
		}

	}
}

func processNewUsers() {
	for user := range newUsersCh {
		processNewUser(user, false)
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

	gui, err = ui.DeployGUI(shutdown, myMsgCh)
	if err != nil {
		log.Fatal(err)
	}

	if *masterMode {
		// Create a new server with the desired parameters.
		srv = server.New("/", *httpPort, msgCh)
		// Start the server (initialize the TCP listener).
		err = srv.Start()
		if err != nil {
			log.Fatal(err)
		}

		// Create a new TCP listener.
		l, err = server.NewListener(*tcpPort, processNewUser, *userName)
		if err != nil {
			log.Fatal(err)
		}
		// Start the listener.
		l.Start()
	} else {
		c, err = client.NewClient(*userName, *masterHost, *tcpPort, msgCh, newUsersCh)
		if err != nil {
			log.Fatal(err)
		}
		initUsers := c.ListenToMasterCommands()
		for _, initUser := range initUsers {
			if initUser != *userName {
				processNewUser(initUser, false)
			}
		}
		go processNewUsers()
	}

	// Writing the incoming data to the "msgCh" channel.
	go processClientsMessages()

	processNewUser(*userName, false)

	go processMyMessages()

	select {}
}
