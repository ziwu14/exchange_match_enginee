package main

import (
	"app/TCPserver"
	"app/command"
	"app/redis"
	"app/xmlParser"

	"sync"
	"fmt"
	"net"
	"time"
	// "runtime"
)

func main() {
	// redis pool
	// runtime.GOMAXPROCS(4)
	redisPool := redis.NewRConnectionPool(
		redis.Config{
			Server:              "redis:6379",
			Password:            "",
			MaxIdle:             100,
			MaxActive:           12000,
			IdleTimeout:         240 * time.Second,
			KEY_PREFIX:          "",
			KEY_DELIMITER:       "",
			KEY_VAR_PLACEHOLDER: "",
		},
	)

	_ = redisPool

	connection := redisPool.Get()
	conn := &connection
	redis.FlushAll(conn)
	connection.Close()

	var readWriteLock sync.RWMutex

	// TCPserver
	server := TCPserver.NewTCPServer(":12345") // set to current ip address (not localhost address)

	server.OnConnectionOpenedCallback = func() {
		fmt.Println("Connection starts")
	}

	server.OnConnectionRecievedNewRequestCallback = func(conn *net.Conn, request []byte) {
		xmlParser := xmlParser.XmlParser{}

		request_in_string := string(request)

		commandList, err := xmlParser.Parse(request_in_string)
		if err != nil {
			(*conn).Write([]byte(fmt.Sprintf("%s", err)))
			return
		}

		commandExecutor := command.CommandListExecutor{Pool: redisPool, ReadWriteLock: &readWriteLock}
		commandExecutor.Execute(commandList)
		response_in_string := commandExecutor.GetResponse()

		(*conn).Write([]byte(response_in_string))
	}

	server.OnConnectionClosedCallback = func(err error) {
		fmt.Println("Connection closed")
		if err != nil {
			fmt.Println("err: ", err)
		}
	}

	server.Listen()
}
