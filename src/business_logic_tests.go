package main

import (
	test "businessLogic"
	"fmt"
	redis "redis"
	"time"

	keyGenerator "uniqueKeyGenerator"

	redigo "github.com/gomodule/redigo/redis"
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

	for i := 0; i < 20; i++ {
		key, _ := keyGenerator.GetNewOrderId(pool)
		fmt.Println(key)
	}

	// 	buyer := "12345"
	// 	seller := "34567"
	// 	createAccount(pool, buyer, 10000, false)
	// 	createAccount(pool, seller, 0, false)
	// 	addBitcoin(pool, buyer, 0, false)
	// 	addBitcoin(pool, seller, 100, false)

	// 	setBuyOrder(pool, "1", buyer, 100, 10, false)
	// 	setSellOrder(pool, "11", seller, 100, 10, false)

	// 	setSellOrder(pool, "12", seller, 4, 10, false)
	// 	setSellOrder(pool, "13", seller, 5, 10, false)
	// 	setSellOrder(pool, "14", seller, 6, 10, false)

	// 	setBuyOrder(pool, "2", buyer, 7, 50, false)

	// 	queryAccount(conn, buyer)
	// 	queryAccount(conn, seller)
	// 	fmt.Println("--------------  buy order -------------------------")
	// 	queryOrder(pool, "1")
	// 	queryOrder(pool, "2")
	// 	fmt.Println("--------------  sell order -------------------------")
	// 	queryOrder(pool, "11")
	// 	queryOrder(pool, "12")
	// 	queryOrder(pool, "13")
	// 	queryOrder(pool, "14")

	// 	fmt.Println("--------------  second -------------------------")

	// 	setBuyOrder(pool, "3", buyer, 8, 10, false)
	// 	setBuyOrder(pool, "4", buyer, 9, 10, false)
	// 	setSellOrder(pool, "15", seller, 6, 50, false)

	// 	queryAccount(conn, buyer)
	// 	queryAccount(conn, seller)
	// 	fmt.Println("--------------  buy order -------------------------")
	// 	queryOrder(pool, "3")
	// 	queryOrder(pool, "4")
	// 	fmt.Println("--------------  sell order -------------------------")
	// 	queryOrder(pool, "15")
	// }

	// func createAccount(pool *redigo.Pool, uid string, balance float64, display bool) {
	// 	test.CreateAccount(pool, uid, balance)
	// 	if display {
	// 		fmt.Printf("Create account: \"%s\", balance: %.2f\n", uid, balance)
	// 	}
}

func addBitcoin(pool *redigo.Pool, uid string, amount float64, display bool) {
	test.SetOrAddSymbolPositionToAccount(pool, uid, "bitcoin", amount)
	if display {
		fmt.Printf("Add bitcoin to account: \"%s\", amount: %.2f\n", uid, amount)
	}
}

func setBuyOrder(pool *redigo.Pool, orderId string, uid string, limitPrice float64, amount float64, display bool) {
	test.SetBuyOrder(pool, orderId, uid, "bitcoin", limitPrice, amount)
	if display {
		fmt.Printf("Set buy order: \"%s\" for account: \"%s\", limitPrice: %.2f, amount: %.2f\n", orderId, uid, limitPrice, amount)
	}
}

func setSellOrder(pool *redigo.Pool, orderId string, uid string, limitPrice float64, amount float64, display bool) {
	test.SetSellOrder(pool, orderId, uid, "bitcoin", limitPrice, amount)
	if display {
		fmt.Printf("Set sell order: \"%s\" for account: \"%s\", limitPrice: %.2f, amount: %.2f\n", orderId, uid, limitPrice, amount)
	}
}

func queryAccount(conn *redigo.Conn, uid string) {
	balance, _ := test.GetAccountBalance(conn, uid)
	amount, _ := test.GetSymbolPosition(conn, uid, "bitcoin")
	fmt.Printf("Query account: \"%s\", balance: %.2f, bitcoin amount:%.2f\n", uid, balance, amount)
}

func queryOrder(pool *redigo.Pool, orderId string) {
	op, ex, can, err := test.QueryOrderStatusAndHistory(pool, orderId)
	fmt.Println("        QueryOrderStatusAndHistory:", orderId, err, "\n", "open:", op, "\n", "exec:", ex, "\n", "canc", can)
}
