package client

import (
	"bufio"
	"net"
	"strconv"
)

// Client contains TCP client's data.
type Client struct {
	conn   *net.TCPConn
	reader *bufio.Reader
}

// NewClient returns new TCP client connection.
func NewClient(userName, host string, port int) (client *Client, err error) {
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
		conn:   conn,
		reader: bufio.NewReader(conn),
	}

	return
}

// Close closes the client connection.
func (c *Client) Close() {
	c.conn.Close()
}
