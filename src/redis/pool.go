package redis

import (
	"log"
	//"os"
	//"os/signal"
	//"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	Pool *redis.Pool
)

type Config struct {
	Server   string
	Password string
	MaxIdle  int // Maximum number of idle connections in the pool.

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait                bool
	KEY_PREFIX          string // prefix to all keys; example is "dev environment name"
	KEY_DELIMITER       string // delimiter to be used while appending keys; example is ":"
	KEY_VAR_PLACEHOLDER string // placeholder to be parsed using given arguments to obtain a final key; example is "?"
}

func NewRConnectionPool(conf Config) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: conf.IdleTimeout,
		MaxActive:   conf.MaxActive,
		Wait:        conf.Wait,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Server)
			if err != nil {
				log.Println("redis-wrapper: Redis: Dial failed", err)
				return nil, err
			}
			if _, err := c.Do("AUTH", conf.Password); err != nil && conf.Password != "" {
				log.Println("redis-wrapper: Redis: AUTH failed", err)
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.Println("redis-wrapper: Unable to ping to redis server:", err)
			}
			return err
		},
	}
}









/*
func init() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis:6379" // It should be redis:6349, instead of :6349
	}
	Pool = newPool(redisHost)
	cleanupHook()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func cleanupHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}
*/