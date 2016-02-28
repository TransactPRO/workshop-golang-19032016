package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/TransactPRO/workshop-golang-19032016/server"
	"github.com/TransactPRO/workshop-golang-19032016/ui"
	"github.com/TransactPRO/workshop-golang-19032016/util"
)

var (
	masterMode = flag.Bool("m", false, "run in master mode")
	masterHost = flag.String("master-host", "127.0.0.1", "master host")
	masterPort = flag.Int("master-port", 8081, "master HTTP port")
)

var (
	gui   *ui.UI
	srv   *server.Server
	msgCh = make(chan util.Message)
)

// shutdown stops all the running services and terminates the process.
func shutdown() {
	log.Println("Shutting down..")

	if srv != nil {
		srv.Stop()
	}

	os.Exit(0)
}

func processIncomingMessages() {
	for msg := range msgCh {
		gui.WriteToChat(fmt.Sprintf("[%s] %s: %s", util.ParseTime(msg.Timestamp), msg.User, msg.Contents))
	}
}

func main() {
	flag.Parse()

	if *masterMode {
		// Create a new server with the desired parameters.
		srv = server.New("/", *masterPort, msgCh)
		// Start the server (initialize the TCP listener).
		srv.Start()
		// The server writes the incoming data to the "msgCh" channel.
		go processIncomingMessages()
	}

	var guiErr error
	gui, guiErr = ui.DeployGUI(shutdown)
	if guiErr != nil {
		log.Fatal(guiErr)
	}

	select {}
}
