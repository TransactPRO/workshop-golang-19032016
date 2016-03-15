package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TransactPRO/workshop-golang-19032016/server"
	"github.com/TransactPRO/workshop-golang-19032016/ui"
	"github.com/TransactPRO/workshop-golang-19032016/util"
)

var (
	userName   = flag.String("u", "", "master host")
	masterMode = flag.Bool("m", false, "run in master mode")
	httpPort   = flag.Int("http-port", 8081, "master's HTTP port")
)

var (
	gui             *ui.UI
	textBoxMessages = make(chan string)
	srv             *server.Server
	clientsMessages = make(chan util.Message)
)

// shutdown stops all the running services and terminates the process.
func shutdown() {
	log.Println("Shutting down..")

	if srv != nil {
		srv.Stop()
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
	} else {

	}

	go processMyMessages()

	select {}
}
