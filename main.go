package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	gui   *ui.UI
	srv   *server.Server
	l     *server.Listener
	c     *client.Client
	msgCh = make(chan util.Message)
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

func processIncomingMessages() {
	for msg := range msgCh {
		gui.WriteToView(ui.ChatView, fmt.Sprintf("[%s] %s: %s", util.ParseTime(msg.Timestamp), msg.User, msg.Contents))
	}
}

func main() {
	flag.Parse()

	if *userName == "" {
		log.Fatal("provide the username!")
	}

	var err error

	if *masterMode {
		// Create a new server with the desired parameters.
		srv = server.New("/", *httpPort, msgCh)
		// Start the server (initialize the TCP listener).
		err = srv.Start()
		if err != nil {
			log.Fatal(err)
		}
		// The server writes the incoming data to the "msgCh" channel.
		go processIncomingMessages()

		// Create a new TCP listener.
		l, err = server.NewListener(*tcpPort, func(newUser string) {
			gui.WriteToView(ui.UsersView, newUser)
		})
		if err != nil {
			log.Fatal(err)
		}
		// Start the listener.
		l.Start()
	} else {
		c, err = client.NewClient(*userName, *masterHost, *tcpPort)
		if err != nil {
			log.Fatal(err)
		}
	}

	gui, err = ui.DeployGUI(shutdown)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
