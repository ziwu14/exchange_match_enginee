package businessLogic

import (
	"fmt"
	redis "app/redis"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
)

const (
	DB_ACCOUNT_PREFIX                   = "account:"
	DB_ACCOUNT_FIELD_BALANCE            = "balance"
	DB_SYMBOL_POSITION_FIELD_AMOUNT     = "amount"
	DB_ORDER_PREFIX                     = "order:"
	DB_ORDER_FIELD_ACCOUNT              = "account"
	DB_ORDER_FIELD_SYMBOL               = "symbol"
	DB_ORDER_FIELD_LIMIT_PRICE          = "limit"
	DB_ORDER_FIELD_ORDER_CURRENT_AMOUNT = "amount"
	DB_ORDER_FIELD_ORDER_INITIAL_AMOUNT = "origAmount"
	DB_ORDER_FIELD_ORDER_TYPE           = "orderType"
	DB_BUY_ORDER_BOOK_PREFIX            = "openBuyOrderBook:"
	DB_SELL_ORDER_BOOK_PREFIX           = "openSellOrderBook:"
	DB_CANCEL_HISTORY_PREFIX            = "order-cancel:"
	DB_CANCEL_HISTORY_FIELD_AMOUNT      = "amount"
	DB_CANCEL_HISOTRY_FIELD_TIME        = "time"
	DB_EXECUTED_HISTORY_PREFIX          = "order-executed:"
	DB_EXECUTED_HISTORY_FIELD_AMOUNT    = "amount"
	DB_EXECUTED_HISTORY_FIELD_LIMIT     = "limit"
	DB_EXECUTED_HISOTRY_FIELD_TIME      = "time"
)

/*
		Create an Account with uid and balance. This function will NOT check if the account exists, User has to MAKE SURE that the account exists.
	input --
		uid: user id, no restriction on the length and characters
		balance: initial user balance, no restriction on the amount, can be negative
*/
func createAccount(conn *redigo.Conn, uid string, balance float64) error {
	return redis.HMSet(conn, DB_ACCOUNT_PREFIX+uid, map[string]interface{}{DB_ACCOUNT_FIELD_BALANCE: balance})
}

/*
		Check an Account with exists.
	input --
		uid: user id, no restriction on the length and characters
*/
func checkAccountExists(conn *redigo.Conn, uid string) (bool, error) {
	return redis.Exists(conn, DB_ACCOUNT_PREFIX+uid)
}

/*
		Get an Account's balance using uid. This function will NOT check if the account exists, User has to MAKE SURE that the account exists.
	input --
		uid: user id, no restriction on the length and characters
	output --
		return the balance in float64
	err --
		from HGET, from strconv.ParseFloat
*/
func GetAccountBalance(conn *redigo.Conn, uid string) (float64, error) {
	balance_in_string, err := redis.HGet(conn, DB_ACCOUNT_PREFIX+uid, DB_ACCOUNT_FIELD_BALANCE)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(balance_in_string, 64)
}

/*
		Increase an Account's balance using uid. This function will NOT check if the account exists, User has to MAKE SURE that the account exists.
	input --
		uid: user id, no restriction on the length and characters
		amount: the amount you want to increase, will accept negative
	output --
		return the balance after increasement in float64
	err --
		from HIncrByFloat, from strconv.ParseFloat
*/
func increaseAccountBalance(conn *redigo.Conn, uid string, amount float64) (float64, error) {
	balance_after_incr_in_string, err := redis.HIncrByFloat(conn, DB_ACCOUNT_PREFIX+uid, DB_ACCOUNT_FIELD_BALANCE, amount)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(balance_after_incr_in_string, 64)
}

/*
		Decrease an Account's balance using uid. This function will NOT check if the account exists, User has to MAKE SURE that the account exists.
	input --
		uid: user id, no restriction on the length and characters
		amount: the amount you want to decrease, will accept negative
	output --
		return the balance after decreasement in float64
	err --
		from HIncrByFloat, from strconv.ParseFloat
*/
func decreaseAccountBalance(conn *redigo.Conn, uid string, amount float64) (float64, error) {
	minus_amount := -amount
	balance_after_incr_in_string, err := redis.HIncrByFloat(conn, DB_ACCOUNT_PREFIX+uid, DB_ACCOUNT_FIELD_BALANCE, minus_amount)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(balance_after_incr_in_string, 64)
}

/*
		Set a symbol position(amount) to an account using uid and symbol name.
		This function will NOT check if the account exists, or the position User has to MAKE SURE that the account exists.
		If the symbol position to that account does not exist, this function will create and set that symbol position to input amount
		If the symbol position already exists, this fucntion will UPDATE that symbol position to input amount
	input --
		uid: user id, no restriction on the length and characters
		symbolName: symbol Name, no restriction on the length and characters
		amount: initial symbol position amount to the account, no restriction on the amount, can be negative
*/
func setSymbolPosition(conn *redigo.Conn, uid string, symbolName string, amount float64) error {
	key := DB_ACCOUNT_PREFIX + uid + ":" + symbolName
	return redis.HMSet(conn, key, map[string]interface{}{DB_SYMBOL_POSITION_FIELD_AMOUNT: amount})
}

/*
		Check an symbol Position with symbolName exists under the account.
	input --
		uid: user id, no restriction on the length and characters
		symbolName: symbol Name, no restriction on the length and characters
*/
func checkSymbolPositionExists(conn *redigo.Conn, uid string, symbolName string) (bool, error) {
	key := DB_ACCOUNT_PREFIX + uid + ":" + symbolName
	return redis.Exists(conn, key)
}

/*
		Get a symbol position amount associated with an account.
	input --
		uid: user id, no restriction on the length and characters
		symbolName: symbol Name, no restriction on the length and characters
	output --
		return the amount in float64
	err --
		from HGET, from strconv.ParseFloat

*/
func GetSymbolPosition(conn *redigo.Conn, uid string, symbolName string) (float64, error) {
	key := DB_ACCOUNT_PREFIX + uid + ":" + symbolName
	amount_in_string, err := redis.HGet(conn, key, DB_SYMBOL_POSITION_FIELD_AMOUNT)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(amount_in_string, 64)
}

/*
		Increase a symbol position amount associated with an account.
		This function will NOT check if the account or the symbol position in the account exists, User has to MAKE SURE that they exist.
	input --
		uid: user id, no restriction on the length and characters
		symbolName: symbol Name, no restriction on the length and characters
		amount: the amount you want to increase, will accept negative
	output --
		return the amount after increasement in float64
	err --
		from HIncrByFloat, from strconv.ParseFloat
*/
func increaseSymbolPosition(conn *redigo.Conn, uid string, symbolName string, amount float64) (float64, error) {
	key := DB_ACCOUNT_PREFIX + uid + ":" + symbolName
	amount_after_incr_in_string, err := redis.HIncrByFloat(conn, key, DB_SYMBOL_POSITION_FIELD_AMOUNT, amount)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(amount_after_incr_in_string, 64)
}

/*
		Decrease a symbol position amount associated with an account.
		This function will NOT check if the account or the symbol position in the account exists, User has to MAKE SURE that they exist.
	input --
		uid: user id, no restriction on the length and characters
		symbolName: symbol Name, no restriction on the length and characters
		amount: the amount you want to decrease, will accept negative
	output --
		return the amount after decreasement in float64
	err --
		from HIncrByFloat, from strconv.ParseFloat
*/
func decreaseSymbolPosition(conn *redigo.Conn, uid string, symbolName string, amount float64) (float64, error) {
	minus_amount := -amount
	key := DB_ACCOUNT_PREFIX + uid + ":" + symbolName
	amount_after_decr_in_string, err := redis.HIncrByFloat(conn, key, DB_SYMBOL_POSITION_FIELD_AMOUNT, minus_amount)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(amount_after_decr_in_string, 64)
}

/*
		Create an buy Order.
		This function will NOT validate anything.(old order, account, symbol position, balance...)
		WARN: If an order with the same orderId exists, the old order will be UPDATED.
	input --
		orderId: order id, no restriction on the length and characters, MAKE SURE it is unique
		uid: user id, no restriction on the length and characters
		symbolName: symbol name, no restriction on the length and characters
		limitPrice: limit price, can be negative
		orderAmount: the symbol position amount you want to buy
*/
func createBuyOrder(conn *redigo.Conn, orderId string, uid string, symbolName string, limitPrice float64, orderAmount float64) error {
	return redis.HMSet(conn,
		DB_ORDER_PREFIX+orderId,
		map[string]interface{}{
			DB_ORDER_FIELD_ACCOUNT:              uid,
			DB_ORDER_FIELD_SYMBOL:               symbolName,
			DB_ORDER_FIELD_LIMIT_PRICE:          limitPrice,
			DB_ORDER_FIELD_ORDER_CURRENT_AMOUNT: orderAmount,
			DB_ORDER_FIELD_ORDER_INITIAL_AMOUNT: orderAmount,
			DB_ORDER_FIELD_ORDER_TYPE:           "buy"})
}

/*
		Create an sell Order.
		This function will NOT validate anything.(old order, account, symbol position, balance...)
		WARN: If an order with the same orderId exists, the old order will be UPDATED.
	input --
		orderId: order id, no restriction on the length and characters, MAKE SURE it is unique
		uid: user id, no restriction on the length and characters
		symbolName: symbol name, no restriction on the length and characters
		limitPrice: limit price, can be negative
		orderAmount: the symbol position amount you want to sell
*/
func createSellOrder(conn *redigo.Conn, orderId string, uid string, symbolName string, limitPrice float64, orderAmount float64) error {
	return redis.HMSet(conn,
		DB_ORDER_PREFIX+orderId,
		map[string]interface{}{
			DB_ORDER_FIELD_ACCOUNT:              uid,
			DB_ORDER_FIELD_SYMBOL:               symbolName,
			DB_ORDER_FIELD_LIMIT_PRICE:          limitPrice,
			DB_ORDER_FIELD_ORDER_CURRENT_AMOUNT: orderAmount,
			DB_ORDER_FIELD_ORDER_INITIAL_AMOUNT: orderAmount,
			DB_ORDER_FIELD_ORDER_TYPE:           "sell"})
}

/*
		Get current amount of an order.
		This function will NOT validate if the orderId exists or not.
	input --
		orderId: order id, no restriction on the length and characters, MAKE SURE it exists
	output --
		return the amount in float64
	err --
		from HGET, from strconv.ParseFloat

*/
func GetOrderAmount(conn *redigo.Conn, orderId string) (float64, error) {
	amount_in_string, err := redis.HGet(conn, DB_ORDER_PREFIX+orderId, DB_ORDER_FIELD_ORDER_CURRENT_AMOUNT)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(amount_in_string, 64)
}

/*
		Get symbolName and orderType of an order.
		This function will NOT validate if the orderId exists or not.
	input --
		orderId: order id, no restriction on the length and characters, MAKE SURE it exists
	output --
		return the symbolName and orderType in string
	err --
		from HMGet

*/
func GetSymbolNameAndOrderType(conn *redigo.Conn, orderId string) ([]string, error) {
	symbolName_n_orderType, err := redis.HMGet(conn, DB_ORDER_PREFIX+orderId, []string{DB_ORDER_FIELD_SYMBOL, DB_ORDER_FIELD_ORDER_TYPE})
	if err != nil {
		return []string{}, err
	}

	return symbolName_n_orderType, nil
}

/*
		Get uid (account id) of an order.
		This function will NOT validate if the orderId exists or not.
	input --
		orderId: order id, no restriction on the length and characters, MAKE SURE it exists
	output --
		return the uid in string
	err --
		from HGet

*/
func GetOrderUid(conn *redigo.Conn, orderId string) (string, error) {
	return redis.HGet(conn, DB_ORDER_PREFIX+orderId, DB_ORDER_FIELD_ACCOUNT)
}

/*
		Get limit price of an order.
		This function will NOT validate if the orderId exists or not.
	input --
		orderId: order id, no restriction on the length and characters, MAKE SURE it exists
	output --
		return the limit price in float64
	err --
		from HGet, strconv.ParseFloat

*/
func GetOrderLimitPrice(conn *redigo.Conn, orderId string) (float64, error) {
	limitPrice_in_string, err := redis.HGet(conn, DB_ORDER_PREFIX+orderId, DB_ORDER_FIELD_LIMIT_PRICE)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(limitPrice_in_string, 64)
}

/*
		Decrease the current amount of an order associated with the orderId.
		This function will NOT check if the order exists, User has to MAKE SURE that it exist.
	input --
		orderId: order id, no restriction on the length and characters, MAKE SURE it exists
		amount: the amount you want to decrease, will accept negative
	output --
		return the current order amount after decreasement in float64
	err --
		from HIncrByFloat, from strconv.ParseFloat
*/
func decreaseOrderAmount(conn *redigo.Conn, orderId string, amount float64) (float64, error) {
	minus_amount := -amount
	amount_after_decr_in_string, err := redis.HIncrByFloat(conn, DB_ORDER_PREFIX+orderId, DB_ORDER_FIELD_ORDER_CURRENT_AMOUNT, minus_amount)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(amount_after_decr_in_string, 64)
}

/*
		Remove an order associated with the orderId.
	input --
		orderId: order id, no restriction on the length and characters
	err --
		from Delete
*/
func removeOrder(conn *redigo.Conn, orderId string) error {
	return redis.Delete(conn, DB_ORDER_PREFIX+orderId)
}

/*
		Check an order with orderId exists.
	input --
		orderId: order id, no restriction on the length and characters
	err --
		from exists
*/
func checkOrderExists(conn *redigo.Conn, orderId string) (bool, error) {
	return redis.Exists(conn, DB_ORDER_PREFIX+orderId)
}

/*
		Add an order reference to buy order book associated with symbolName.
		This function will not check the existence of the order.
		If orderId exists, the old order reference will be UPDATED.
	input --
		symbolName: the symbol that this order belongs to.
		orderId: order id, no restriction on the length and characters
		limitPrice: the limit price of this order
*/
func AddBuyOrderToBuyOrderBook(conn *redigo.Conn, symbolName string, orderId string, limitPrice float64) error {
	return redis.ZAdd(conn, DB_BUY_ORDER_BOOK_PREFIX+symbolName, limitPrice, orderId)
}

/*
		Add an order reference to sell order book associated with symbolName.
		This function will not check the existence of the order.
		If orderId exists, the old order reference will be UPDATED.
	input --
		symbolName: the symbol that this order belongs to.
		orderId: order id, no restriction on the length and characters
		limitPrice: the limit price of this order
*/
func AddSellOrderToSellOrderBook(conn *redigo.Conn, symbolName string, orderId string, limitPrice float64) error {
	return redis.ZAdd(conn, DB_SELL_ORDER_BOOK_PREFIX+symbolName, limitPrice, orderId)
}

/*
		Remove an order reference from a buy order book associated with symbolName.
		This function will not check the existence of the order.
	input --
		symbolName: the symbol that this order belongs to.
		orderId: order id, no restriction on the length and characters
*/
func removeBuyOrderFromBuyOrderBook(conn *redigo.Conn, symbolName string, orderId string) error {
	return redis.ZRem(conn, DB_BUY_ORDER_BOOK_PREFIX+symbolName, orderId)
}

/*
		Remove an order reference from a sell order book associated with symbolName.
		This function will not check the existence of the order.
	input --
		symbolName: the symbol that this order belongs to.
		orderId: order id, no restriction on the length and characters
*/
func removeSellOrderFromSellOrderBook(conn *redigo.Conn, symbolName string, orderId string) error {
	return redis.ZRem(conn, DB_SELL_ORDER_BOOK_PREFIX+symbolName, orderId)
}

/*
		Return the orderId with maximum limit price and its limitPrice in a buy order book associated with symbolName.
		This function will not check the existence of the order book.
	input --
		symbolName: the buy order book's symbol that you want to peek
*/
func peekBuyOrderWithMaxPriceInBuyOrdrerBook(conn *redigo.Conn, symbolName string) (string, float64, error) {
	orderId_n_limitPrice, err := redis.ZRevRange(conn, DB_BUY_ORDER_BOOK_PREFIX+symbolName, 0, 0, true)
	if err != nil {
		return "", 0, err
	} else if len(orderId_n_limitPrice) == 0 {
		return "", 0, fmt.Errorf("empty buy order book")
	}

	orderId := orderId_n_limitPrice[0]
	var limitPrice float64
	limitPrice, err = strconv.ParseFloat(orderId_n_limitPrice[1], 64)
	if err != nil {
		return "", 0, err
	}

	return orderId, limitPrice, nil
}

/*
		Return the orderId with minimum limit price and its limitPrice in a sell order book associated with symbolName.
		This function will not check the existence of the order book.
	input --
		symbolName: the sell order book's symbol that you want to peek
*/
func peekSellOrderWithMinPriceInSellOrdrerBook(conn *redigo.Conn, symbolName string) (string, float64, error) {
	orderId_n_limitPrice, err := redis.ZRange(conn, DB_SELL_ORDER_BOOK_PREFIX+symbolName, 0, 0, true)
	if err != nil {
		return "", 0, err
	} else if len(orderId_n_limitPrice) == 0 {
		return "", 0, fmt.Errorf("empty sell order book")
	}

	orderId := orderId_n_limitPrice[0]
	var limitPrice float64
	limitPrice, err = strconv.ParseFloat(orderId_n_limitPrice[1], 64)
	if err != nil {
		return "", 0, err
	}

	return orderId, limitPrice, nil
}

/*
		Check if a buy order book is empty. If the book does not exist, true is returned.
	input --
		symbolName: the buy order book symbol that you want to check
*/
func isBuyOrderBookEmpty(conn *redigo.Conn, symbolName string) (bool, error) {
	number_of_elements, err := redis.ZCard(conn, DB_BUY_ORDER_BOOK_PREFIX+symbolName)
	if err != nil {
		return true, err
	}

	if number_of_elements == 0 {
		return true, nil
	}

	return false, nil
}

/*
		Check if a sell order book is empty. If the book does not exist, true is returned.
	input --
		symbolName: the buy order book symbol that you want to check
*/
func isSellOrderBookEmpty(conn *redigo.Conn, symbolName string) (bool, error) {
	number_of_elements, err := redis.ZCard(conn, DB_SELL_ORDER_BOOK_PREFIX+symbolName)
	if err != nil {
		return true, err
	}

	if number_of_elements == 0 {
		return true, nil
	}

	return false, nil
}

/*
		Insert a cancel order history tuple to cancel order histories.
		This function will NOT check if the history exists.
		If a history with same order id exists, this function will UPDATE the old history.
	input --
		orderId: order id, no restriction on the length and characters, MAKE SURE it is unique in cancel histories
		amount: the order's current order amount before cancellation
		time: cancellation time
*/
func insertCancelledOrderToCancelHistory(conn *redigo.Conn, orderId string, amount float64, time string) error {
	return redis.HMSet(conn,
		DB_CANCEL_HISTORY_PREFIX+orderId,
		map[string]interface{}{
			DB_CANCEL_HISTORY_FIELD_AMOUNT: amount,
			DB_CANCEL_HISOTRY_FIELD_TIME:   time})
}

/*
		Query a cancel order history tuple for its current amount at cancelled time and cancelled time.
		This function will NOT check if the history exists. MAKE SURE that the history EXIST.
	input --
		orderId: order id, no restriction on the length and characters, MAKE SURE it is unique in cancel histories
*/
func getAmountAndTimeForCancelledOrderFromCancelHistory(conn *redigo.Conn, orderId string) (float64, string, error) {
	amount_n_time, err := redis.HMGet(conn, DB_CANCEL_HISTORY_PREFIX+orderId, []string{DB_CANCEL_HISTORY_FIELD_AMOUNT, DB_CANCEL_HISOTRY_FIELD_TIME})
	if err != nil {
		return 0, "", err
	} else if len(amount_n_time) == 0 {
		return 0, "", fmt.Errorf("no cancelled order with this orderId")
	}

	amount_in_string := amount_n_time[0]
	time := amount_n_time[1]

	var amount float64
	amount, err = strconv.ParseFloat(amount_in_string, 64)
	if err != nil {
		return 0, "", err
	}

	return amount, time, nil
}

/*
		Check a cancel order history tuple with orderId exists.
	input --
		orderId: order id, no restriction on the length and characters
*/
func cancelledOrderExists(conn *redigo.Conn, orderId string) (bool, error) {
	exists, err := redis.Exists(conn, DB_CANCEL_HISTORY_PREFIX+orderId)
	if err != nil {
		return false, err
	}

	return exists, nil
}

/*
		WARN: By the time we write this document, we did no know how to write TRANSACTION, so we will not handle errors in consecutive RPUSHs

		Insert an executed order history tuple to executed order histories.
	input --
		orderId: order id, no restriction on the length and characters
		amount: the order's executed order amount
		limitPrice: the order's executed limit price.
		time: executed time
*/
func InsertExcutedOrderToExcutedHistory(conn *redigo.Conn, orderId string, amount float64, limitPrice float64, time string) error {
	amount_in_string := fmt.Sprintf("%f", amount)
	limitPrice_in_string := fmt.Sprintf("%f", limitPrice)
	redis.RPush(conn, DB_EXECUTED_HISTORY_PREFIX+orderId, amount_in_string)
	redis.RPush(conn, DB_EXECUTED_HISTORY_PREFIX+orderId, limitPrice_in_string)
	redis.RPush(conn, DB_EXECUTED_HISTORY_PREFIX+orderId, time)
	return nil
}

/*
		Check a executed order history list with orderId exists.
	input --
		orderId: order id, no restriction on the length and characters
*/
func executedOrderExists(conn *redigo.Conn, orderId string) (bool, error) {
	exists, err := redis.Exists(conn, DB_EXECUTED_HISTORY_PREFIX+orderId)
	if err != nil {
		return false, err
	}

	return exists, nil
}

/*
		Query an executed order history slice list.
		list eg: (EO: executed order)
			amount of EO1 --> limit price of EO1 --> time of EO1 --> amount of EO2 --> limit price of EO2 --> time of EO2 --> ...
		This function will NOT check if the history list exists. MAKE SURE that the history list EXIST.
	input --
		orderId: order id, no restriction on the length and characters.
*/
func GetExecutedOrderSliceList(conn *redigo.Conn, orderId string) ([]string, error) {
	return redis.LRange(conn, DB_EXECUTED_HISTORY_PREFIX+orderId, 0, -1)
}
