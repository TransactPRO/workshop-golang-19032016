package main

import (
	"log"
	"os"

	"github.com/TransactPRO/workshop-golang-19032016/ui"
)

var (
	gui *ui.UI
)

func shutdown() {
	log.Println("Shutting down..")

	os.Exit(0)
}

func main() {
	var guiErr error
	gui, guiErr = ui.DeployGUI(shutdown)
	if guiErr != nil {
		log.Fatal(guiErr)
	}

	select {}
}
