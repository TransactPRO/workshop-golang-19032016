package client

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/TransactPRO/workshop-golang-19032016/util"
)

// Client contains TCP client's data.
type Client struct {
	conn       *net.TCPConn
	reader     *bufio.Reader
	msgCh      chan util.Message
	newUsersCh chan string
}

// NewClient returns new TCP client connection.
func NewClient(userName, host string, port int, msgCh chan util.Message, newUsersCh chan string) (client *Client, err error) {
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp4", host+":"+strconv.Itoa(port))
	if err != nil {
		return
	}

	var conn *net.TCPConn
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}

	conn.Write([]byte(userName + "\n"))

	client = &Client{
		conn:       conn,
		reader:     bufio.NewReader(conn),
		msgCh:      msgCh,
		newUsersCh: newUsersCh,
	}

	return
}

// ListenToMasterCommands receives master commands.
func (c *Client) ListenToMasterCommands() (initUsers []string) {
	initResult, err := c.reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}

	initErrr := json.Unmarshal(initResult, &initUsers)
	if initErrr != nil {
		log.Fatal(initErrr)
	}

	go func() {
		for {
			result, err := c.reader.ReadBytes('\n')

			if err != nil {
				if err == io.EOF {
					log.Println("master has disconnected")
					return
				}
				log.Println("failed to process master command")
				continue
			}

			var cmd util.Command
			unmarshalErr := json.Unmarshal(result, &cmd)
			if unmarshalErr != nil {
				log.Println(unmarshalErr)
				continue
			}

			switch cmd.ID {
			case "MSG":
				c.msgCh <- cmd.Message
			case "USER":
				c.newUsersCh <- cmd.OriginUser
			}
		}
	}()

	return
}

// Close closes the client connection.
func (c *Client) Close() {
	c.conn.Close()
}
