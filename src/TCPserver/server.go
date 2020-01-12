package TCPserver

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Connections handles request from client and send response to client
// ------------------------------------------------------------------------------------------
type Connection struct {
	conn                                   net.Conn
	onConnectionOpenedCallback             func()
	onConnectionRecievedNewRequestCallback func(conn *net.Conn, request []byte)
	onConnectionClosedCallback             func(err error)
}

// connection read data(XML) from client
func (c *Connection) listen() {
	c.onConnectionOpenedCallback()

	reader := bufio.NewReader(c.conn)

	message_length, err := reader.ReadString('\n')
	if err != nil {
		c.conn.Write([]byte("xml wrong format in first line"))
		c.conn.Close()
		c.onConnectionClosedCallback(err)
		return
	}

	message_length = strings.TrimSuffix(message_length, "\n")

	// message_length should indicate number of bytes of XML to read
	len_msg, err := strconv.Atoi(message_length)
	if err != nil {
		c.conn.Write([]byte("xml wrong format in first line"))
		c.conn.Close()
		c.onConnectionClosedCallback(err)
		return
	}

	request := make([]byte, len_msg)

	bytes_read, err := reader.Read(request)
	// ensure that bytes read matches length specified of XML request
	if err != nil || bytes_read != len_msg {
		c.conn.Write([]byte("xml byte indicator and xml size mismatch"))
		c.conn.Close()
		c.onConnectionClosedCallback(err)
		return
	}

	c.onConnectionRecievedNewRequestCallback(&c.conn, request)

	c.conn.Close()
	c.onConnectionClosedCallback(err)
}

// TCP server
// ------------------------------------------------------------------------------------------
type server struct {
	address                                string // Address to open connection
	OnConnectionOpenedCallback             func()
	OnConnectionRecievedNewRequestCallback func(conn *net.Conn, request []byte)
	OnConnectionClosedCallback             func(err error)
}

// Creates a tcp server instance
func NewTCPServer(address string) *server {
	log.Info("Creating server with address: ", address)
	server := &server{
		address: address,
	}

	server.OnConnectionOpenedCallback = func() {}
	server.OnConnectionRecievedNewRequestCallback = func(conn *net.Conn, request []byte) {}
	server.OnConnectionClosedCallback = func(err error) {}

	return server
}

// server listens for new internet connection, create a socket and transfer later work to Connection
func (s *server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client_connection := &Connection{
			conn:                                   conn,
			onConnectionOpenedCallback:             s.OnConnectionOpenedCallback,
			onConnectionRecievedNewRequestCallback: s.OnConnectionRecievedNewRequestCallback,
			onConnectionClosedCallback:             s.OnConnectionClosedCallback,
		}

		go client_connection.listen()
	}
}

// main function
// ------------------------------------------------------------------------------------------
func main() {
	server := NewTCPServer("10.197.182.188:12345") // set to current ip address (not localhost address)

	server.OnConnectionOpenedCallback = func() {
		fmt.Println("Connection starts")
	}

	server.OnConnectionRecievedNewRequestCallback = func(conn *net.Conn, request []byte) {
		fmt.Printf("Connection received request from client: %s", string(request))
		(*conn).Write(request)
	}

	server.OnConnectionClosedCallback = func(err error) {
		fmt.Println("Connection closed")
		if err != nil {
			fmt.Println("err: ", err)
		}
	}

	server.Listen()
}
