package uniqueKeyGenerator

import (
	redis "app/redis"

	redigo "github.com/gomodule/redigo/redis"
)

const (
	DB_KEY_FOR_ORDER_ID_GENERATOR = "orderIdCounter"
)

func GetNewOrderId(pool *redigo.Pool) (int, error) {
	connection := pool.Get()
	defer connection.Close()
	conn := (&connection)

	return redis.Incr(conn, DB_KEY_FOR_ORDER_ID_GENERATOR)
}
