package main

import (
	test "businessLogic"
	"fmt"
	redis "redis"
	"time"
)

func main() {
	pool := redis.NewRConnectionPool(
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

	connection := pool.Get()
	defer connection.Close()
	conn := (&connection)

	redis.FlushAll(conn)

	for i := 0; i < 100000; i++ {
		test.InsertExcutedOrderToExcutedHistory(conn, "3", 100, 100, "time")
	}

	exists, _ := test.ExecutedOrderExists(conn, "3")
	fmt.Println("InsertExcutedOrderToExcutedHistory & ExecutedOrderExists <true>:", exists)

	exists, _ = test.ExecutedOrderExists(conn, "2")
	fmt.Println("ExecutedOrderExists <false>:", exists)

	historyList, _ := test.GetExecutedOrderSliceList(conn, "3")
	fmt.Println("length of hitoryList <9000>", len(historyList))
	listLength, _ := redis.LLen(conn, "order-executed:"+"3")
	fmt.Println("length of list in redis <9000>", listLength)
}
