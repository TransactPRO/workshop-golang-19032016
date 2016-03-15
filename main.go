package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TransactPRO/workshop-golang-19032016/ui"
	"github.com/TransactPRO/workshop-golang-19032016/util"
)

var (
	userName = flag.String("u", "", "master host")
)

var (
	gui             *ui.UI
	textBoxMessages = make(chan string)
)

// shutdown stops all the running services and terminates the process.
func shutdown() {
	log.Println("Shutting down..")

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

	go processMyMessages()

	select {}
}
