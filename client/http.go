package client

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TransactPRO/workshop-golang-19032016/util"
)

// SendMessageToMaster sends the message to master.
func SendMessageToMaster(msg util.Message, host string, port int) {
	url := "http://" + host + ":" + strconv.Itoa(port)

	b, marshalErr := json.Marshal(msg)
	if marshalErr != nil {
		log.Println(marshalErr)
	}

	req, reqErr := http.NewRequest("POST", url, bytes.NewReader(b))
	if reqErr != nil {
		log.Println(reqErr)
	}

	client := new(http.Client)
	_, doErr := client.Do(req)
	if doErr != nil {
		log.Println(doErr)
	}
}
