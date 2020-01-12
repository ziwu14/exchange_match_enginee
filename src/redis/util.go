package redis

import (
	"fmt"

	redis "github.com/gomodule/redigo/redis"
)

func Ping(conn *redis.Conn) error {
	_, err := redis.String((*conn).Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}

// Create a Hash with given key, set a field with corresponding value.
// If the hash exists, set/add a field/value pair to it.
// workon redis dataType: Hash
func HSet(conn *redis.Conn, key string, field string, value interface{}) error {
	_, err := (*conn).Do("HSET", key, field, value)
	return err
}

// Create a Hash with given key, set fields with corresponding values
// If the hash exists, set/add field/value pairs to it.
// workon redis dataType: Hash
func HMSet(conn *redis.Conn, key string, fieldValueMap map[string]interface{}) error {
	var args = []interface{}{key}
	for f, v := range fieldValueMap {
		args = append(args, f, v)
	}
	_, err := (*conn).Do("HMSET", args...)
	return err
}

// DEPRECATED BECAUSE RESULT OUT OF ORDER: Get All fields/values from a Hash with the given key
// WARN: HGetAll will return [field/value] in a NON-PREDICTABLE order
// eg for return value: [field1 value1 field2 value2 field3 value3 ...]
// workon redis dataType: Hash
// returnType: []string
func HGetAll(conn *redis.Conn, key string) ([]string, error) {

	data, err := redis.Strings((*conn).Do("HGETALL", key))
	return data, err
}

// Get the value under field from a Hash associated with the given key
// workon redis dataType: Hash
// returnType: string
func HGet(conn *redis.Conn, key string, field string) (string, error) {
	value, err := redis.String((*conn).Do("HGET", key, field))
	return value, err
}

// Get the values under field from a Hash associated with the given key
// Result is returned as the same order as fields input: [field1, field2] --> [value1, value2]
// workon redis dataType: Hash
// returnType: []string
func HMGet(conn *redis.Conn, key string, fields []string) ([]string, error) {
	args := []interface{}{key}
	for _, field := range fields {
		args = append(args, field)
	}
	values, err := redis.Strings((*conn).Do("HMGET", args...))
	return values, err
}

// HExists checks if a hash associated with key and with given field exists
// workon redis dataType: Hash
func HExists(conn *redis.Conn, key string, field string) (bool, error) {
	exists, err := redis.Bool((*conn).Do("HEXISTS", key, field))
	return exists, err
}

// HIncrByFloat increases a field(should be float64 field) by given amount
// return field value after increasement in string
// workon redis dataType: Hash
func HIncrByFloat(conn *redis.Conn, key string, field string, amount float64) (string, error) {
	amountAfterIncrease, err := redis.String((*conn).Do("HINCRBYFLOAT", key, field, amount))
	return amountAfterIncrease, err
}

// Incr return the current value and increase the value by 1
// If the key does not exist, the key will be set to 0 first, and then incr is called TWICE
// So the value returned for calling Incr for the first time will be 1
// MAKE SURE that the value can be represented as integer
// workon redis dataType: Hash
func Incr(conn *redis.Conn, key string) (int, error) {
	return redis.Int((*conn).Do("INCR", key))
}

// Exists checks if a hash associated with key exists
// workon redis dataType: Hash
func Exists(conn *redis.Conn, key string) (bool, error) {
	ok, err := redis.Bool((*conn).Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

// Delete remove a hash associated with key
// workon redis dataType: Hash
func Delete(conn *redis.Conn, key string) error {
	_, err := (*conn).Do("DEL", key)
	return err
}

// ZAdd adds key with given score(normally float, or a string which can parsed to a float) to a sorted set named setName
// If set named setName does not exist, ZAdd creates a set named setName
// If key exists, ZAdd updates the score to given value
// workon redis dataType: Sorted Set
func ZAdd(conn *redis.Conn, setName string, score interface{}, key string) error {
	_, err := (*conn).Do("ZADD", setName, score, key)
	return err
}

// ZRem remove a key from a Sorted set named setName
// workon redis dataType: Sorted Set
func ZRem(conn *redis.Conn, setName string, key string) error {
	_, err := (*conn).Do("ZREM", setName, key)
	return err
}

// ZRange retrieves a list of keys from a Sorted set named setName, from start to stop
// Order of results: from lowest to highest
// -1 == last element
// 0 == first element
// eg: retrieve first 5 elements (low to high): start=0, stop=4
// If start > stop, returns an empty list
// If stop > last element rank, treat stop as the last element rank
// If the set is empty, an EMPTY []string is returned
// workon redis dataType: Sorted Set
func ZRange(conn *redis.Conn, setName string, start int, stop int, withScores bool) ([]string, error) {
	if withScores {
		return redis.Strings((*conn).Do("ZRANGE", setName, start, stop, "WITHSCORES"))
	}

	return redis.Strings((*conn).Do("ZRANGE", setName, start, stop))
}

// ZRevRange retrieves a list of keys from a Sorted set named setName, from start to stop
// Order of results: from highest to lowest
// -1 == last element
// 0 == first element
// eg: retrieve first 5 elements (high to low): start=0, stop=4
// If start > stop, returns an empty list
// If stop > last element rank, treat stop as the last element rank
// If the set is empty, an EMPTY []string is returned
// workon redis dataType: Sorted Set
func ZRevRange(conn *redis.Conn, setName string, start int, stop int, withScores bool) ([]string, error) {
	if withScores {
		return redis.Strings((*conn).Do("ZREVRANGE", setName, start, stop, "WITHSCORES"))
	}

	return redis.Strings((*conn).Do("ZREVRANGE", setName, start, stop))
}

/*
	Zcard returns the number of elements of the sorted set.
	If the set does not exist, 0 is returned.
	workon redis dataType: Sorted Set
*/
func ZCard(conn *redis.Conn, setName string) (int, error) {
	return redis.Int((*conn).Do("ZCARD", setName))
}

// LPush pre-appends a node to the head of a list named listName
// If list named listName does not exists, it will be created
// workon redis dataType: List
func LPush(conn *redis.Conn, listName string, node string) error {
	_, err := (*conn).Do("LPUSH", listName, node)
	return err
}

// RPush appends a node to the tail of a list named listName
// If list named listName does not exists, it will be created
// eg: RPush(1, 2, 3) ==>  head -> 1, 2, 3 -> tail
// 	   LRange(0 ~ 2) ==> [1, 2, 3]
// tip: use RPush with LRange
// workon redis dataType: List
func RPush(conn *redis.Conn, listName string, node string) error {
	_, err := (*conn).Do("RPUSH", listName, node)
	return err
}

// LRange query a list named with listName, and returns results as []string
// -1 == last element
// 0 == first element
// If start > stop, returns an empty list
// If stop > last element rank, treat stop as the last element rank
// workon redis dataType: List
func LRange(conn *redis.Conn, listName string, start int, stop int) ([]string, error) {
	nodes, err := redis.Strings((*conn).Do("LRANGE", listName, start, stop))
	return nodes, err
}

// LLen returns a length of the list named with listName
// If list does not exist, 0 is returned
// If the container associated with listName is not a list, err is returned
// workon redis dataType: List
func LLen(conn *redis.Conn, listName string) (int, error) {
	return redis.Int((*conn).Do("LLEN", listName))
}

// FlushAll flushes all data in Redis db
// BE CAREFUL WITH THIS FUNCTION
func FlushAll(conn *redis.Conn) {
	(*conn).Do("FLUSHALL")
	return
}
