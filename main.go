package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/TransactPRO/workshop-golang-19032016/ui"
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

	for msg := range textBoxMessages {
		fmt.Println(msg)
	}

	select {}
}
