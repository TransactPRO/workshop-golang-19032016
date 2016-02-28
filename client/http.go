package client

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TransactPRO/workshop-golang-19032016/util"
)

// SendMessageToMaster sends the message to master.
func (c *Client) SendMessageToMaster(msg util.Message, host string, port int) {
	encodedMsg := util.EncodeMessage(msg.Contents)
	_, err := http.Get("http://" + host + ":" + strconv.Itoa(port) + "/?user=" + msg.User + "&msg=" + encodedMsg)
	if err != nil {
		log.Println(err)
	}
}
