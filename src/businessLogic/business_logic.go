package businessLogic

import (
	"fmt"
	"math"

	redigo "github.com/gomodule/redigo/redis"
)

const (
	ORDER_TYPE_BUY  = "buy"
	ORDER_TYPE_SELL = "sell"
)

type CancelledOrderHistoryTuple struct {
	CancelledAmount string
	CancelledTime   string
}

type ExecutedOrderHistoryTuple struct {
	TransactionAmount string
	TransactionPrice  string
	TransactionTime   string
}

type OpenOrderTuple struct {
	CurrentAmount string
}

/*
		CreateAccount will create an account in redis with uid and balance.
	input --
		uid: user id, a base-10 digit sequence
		balance: should be non-negative float(>= 0)
	output --
		error:
		if uid or balance does not meet with input restriction, an error message will be returned
		if database fails to create the account, an error message will be returned
		if no error returns, an account is successfully created in redis
*/
func CreateAccount(pool *redigo.Pool, uid string, balance float64) error {
	if !isBase10NumberSequense(uid) || balance < 0 {
		return fmt.Errorf("invalid id or balance")
	}

	connection := pool.Get()
	defer connection.Close()
	conn := (&connection)

	accountExists, err := checkAccountExists(conn, uid)
	if err != nil || accountExists {
		return fmt.Errorf("user already exists")
	}

	err = createAccount(conn, uid, balance)
	if err != nil {
		return fmt.Errorf("database error to create an account")
	}

	return nil
}

/*
		SetOrAddSymbolPositionToAccount will set an account's symbol position to the amount. Symbol is specified by symbolName.
		If this account already has symbol position for this symbolName, then the amount will be added to the account's symbol position.
	input --
		uid: user id, a base-10 digit sequence
		symbolName: string
		amount: should be non-negative float(>= 0)
	output --
		error:
		if uid does not exist, an error message will be returned
		if amount does not meet input restriction, an error message will be returned
		if database fails to create the symbol position, an error message will be returned
		if no error returns, the symbol position is successfully created under the account in redis
*/
func SetOrAddSymbolPositionToAccount(pool *redigo.Pool, uid string, symbolName string, amount float64) error {
	connection := pool.Get()
	defer connection.Close()
	conn := (&connection)

	exists, err := checkAccountExists(conn, uid)
	if err != nil || !exists {
		return fmt.Errorf("user doesn't exist")
	}

	if amount < 0 {
		return fmt.Errorf("invalid amount")
	}

	exists, err = checkSymbolPositionExists(conn, uid, symbolName)
	if err != nil {
		return fmt.Errorf("database error to create/add symbol")
	}

	if exists {
		_, err = increaseSymbolPosition(conn, uid, symbolName, amount)
		if err != nil {
			return fmt.Errorf("database error to create/add symbol")
		}
	} else {
		err = setSymbolPosition(conn, uid, symbolName, amount)
		if err != nil {
			return fmt.Errorf("database error to create/add symbol")
		}
	}

	return nil
}

/*
		SetBuyOrder will set a buy order for an account. The buy order is for symbol: symbolName, and is set with limitPrice and amount.
		If created successfully, the account's balance will be deducted by limitPrice * amount
	input --
		orderId: order id, MUST BE UNIQUE, THE UNIQUENESS IS MAINTAINED BY THE USER OF THIS FUNCTION.
				 If an order with the same order id exists, it will be UPDATED
		uid: user id, a base-10 digit sequence
		symbolName: string
		limitPrice: should be non-negative float(> 0)
		amount: should be non-negative float(> 0)
	output --
		error:
		if uid does not exist, an error message will be returned
		if amount or limitPrice does not meet input restriction, an error message will be returned
		if the account's balance is insufficient to create the order, an error message will be returned
		if database fails to create the symbol position, an error message will be returned
		if no error returns, the buy order is successfully created under the account in redis
*/
func SetBuyOrder(pool *redigo.Pool, orderId string, uid string, symbolName string, limitPrice float64, amount float64) error {
	connection := pool.Get()
	defer connection.Close()
	conn := (&connection)

	exists, err := checkAccountExists(conn, uid)
	if err != nil || !exists {
		return fmt.Errorf("user doesn't exist")
	}

	if amount <= 0 || limitPrice <= 0 {
		return fmt.Errorf("invalid amount or limit price")
	}

	var accountBalance float64
	accountBalance, err = GetAccountBalance(conn, uid)
	payment := limitPrice * amount
	if accountBalance < payment {
		return fmt.Errorf("insufficient fund")
	}

	err = createBuyOrder(conn, orderId, uid, symbolName, limitPrice, amount)
	if err != nil {
		return fmt.Errorf("database error to create buy order")
	}
	err = AddBuyOrderToBuyOrderBook(conn, symbolName, orderId, limitPrice)
	if err != nil {
		return fmt.Errorf("database error to add buy order to order book")
	}

	_, err = decreaseAccountBalance(conn, uid, payment)
	if err != nil {
		return fmt.Errorf("database error when deducting balance from account")
	}

	MatchOrder(conn, orderId, uid, symbolName, limitPrice, amount, "buy")

	return nil
}

/*
		SetSellOrder will set a sell order for an account. The sell order is for symbol: symbolName, and is set with limitPrice and amount.
		If created successfully, the account's symbol position for this symbol will be deducted by amount.
	input --
		orderId: order id, MUST BE UNIQUE, THE UNIQUENESS IS MAINTAINED BY THE USER OF THIS FUNCTION.
				 If an order with the same order id exists, it will be UPDATED
		uid: user id, a base-10 digit sequence
		symbolName: string
		limitPrice: should be non-negative float(> 0)
		amount: should be non-negative float(> 0)
	output --
		error:
		if uid does not exist, an error message will be returned
		if amount or limitPrice does not meet input restriction, an error message will be returned
		if the account's symbol position for this symbol is insufficient to create the order, an error message will be returned
		if database fails to create the symbol position, an error message will be returned
		if no error returns, the buy order is successfully created under the account in redis
*/
func SetSellOrder(pool *redigo.Pool, orderId string, uid string, symbolName string, limitPrice float64, amount float64) error {
	connection := pool.Get()
	defer connection.Close()
	conn := (&connection)

	exists, err := checkAccountExists(conn, uid)
	if err != nil || !exists {
		return fmt.Errorf("user doesn't exist")
	}

	exists, err = checkSymbolPositionExists(conn, uid, symbolName)
	if err != nil || !exists {
		return fmt.Errorf("symbol position doesn't exist under this account")
	}

	if amount <= 0 || limitPrice <= 0 {
		return fmt.Errorf("invalid amount or limit price")
	}

	var symbolPositionInAccount float64
	symbolPositionInAccount, err = GetSymbolPosition(conn, uid, symbolName)
	if err != nil || symbolPositionInAccount < amount {
		return fmt.Errorf("insufficient symbols")
	}

	err = createSellOrder(conn, orderId, uid, symbolName, limitPrice, amount)
	if err != nil {
		return fmt.Errorf("database error to create sell order")
	}
	err = AddSellOrderToSellOrderBook(conn, symbolName, orderId, limitPrice)
	if err != nil {
		return fmt.Errorf("database error to add sell order to order book")
	}

	_, err = decreaseSymbolPosition(conn, uid, symbolName, amount)
	if err != nil {
		return fmt.Errorf("database error when deducting amount from symbol")
	}

	MatchOrder(conn, orderId, uid, symbolName, limitPrice, amount, "sell")

	return nil
}

/*
		CancelOpenOrder cancels an open order.
	input --
		orderId: order id.
	output --
		err:
		If no open order with order id exists, an error message is returned
		if fails to retrieve order, remove open order, or insert order history from database, an error message will be returned
*/
func CancelOpenOrder(pool *redigo.Pool, orderId string) error {
	connection := pool.Get()
	defer connection.Close()
	conn := (&connection)
	exists, err := checkOrderExists(conn, orderId)
	if err != nil || !exists {
		return fmt.Errorf("open order with this order id does not exist")
	}



	var symbolName_n_orderType []string
	symbolName_n_orderType, err = GetSymbolNameAndOrderType(conn, orderId)
	if err != nil {
		return fmt.Errorf("database error when retrieving symbol name and order type")
	}

	var amount float64
	amount, err = GetOrderAmount(conn, orderId)
	if err != nil {
		return fmt.Errorf("database error when getting order amount")
	}
	var price float64
	price, err = GetOrderLimitPrice(conn, orderId)
	if err != nil {
		return fmt.Errorf("database error when getting order price")
	}
	var uid string
	uid, err = GetOrderUid(conn, orderId)
	if err != nil {
		return fmt.Errorf("database error when getting order uid")
	}

	symbolName := symbolName_n_orderType[0]
	orderType := symbolName_n_orderType[1]


	if orderType == ORDER_TYPE_BUY {
		_, err = increaseAccountBalance(conn, uid, price * amount)
		if err != nil {
			return fmt.Errorf("database error when return money to buyer")
		}
		err = removeBuyOrderFromBuyOrderBook(conn, symbolName, orderId)
		if err != nil {
			return fmt.Errorf("database error when removing buy order from buy order book")
		}
	} else {
		_, err = increaseSymbolPosition(conn, uid, symbolName, amount)
		if err != nil {
			return fmt.Errorf("database error when return symbol to seller")
		}
		err = removeSellOrderFromSellOrderBook(conn, symbolName, orderId)
		if err != nil {
			return fmt.Errorf("database error when removing sell order from sell order book")
		}
	}

	err = removeOrder(conn, orderId)
	if err != nil {
		return fmt.Errorf("database error when removing order from orders")
	}

	currentTimeInString := getCurrentTimeInString()

	err = insertCancelledOrderToCancelHistory(conn, orderId, amount, currentTimeInString)
	if err != nil {
		return fmt.Errorf("database error when inserting cancelled order to cancalled order history")
	}

	return nil
}

/*
		QueryOrderStatusAndHistory query open order, executed history and cancelled order history with order id.
	input --
		orderId: order id.
	output --
		a list of open order tuples, a list of executed order history tuples, a list of cancelled order history tuples,
		eg: if no open order is found for this order id(i.e. the order has been cancelled), the list of open order tuples will be empty
		err:
		If no open order with order id exists, an error message is returned
*/
func QueryOrderStatusAndHistory(pool *redigo.Pool, orderId string) ([]OpenOrderTuple, []ExecutedOrderHistoryTuple, []CancelledOrderHistoryTuple, error) {
	connection := pool.Get()
	defer connection.Close()
	conn := (&connection)

	var exists_in_executed_history, exists_in_cancel_history, exists_in_open_orders bool
	var err error
	exists_in_executed_history, err = executedOrderExists(conn, orderId)
	if err != nil {
		return []OpenOrderTuple{}, []ExecutedOrderHistoryTuple{}, []CancelledOrderHistoryTuple{}, fmt.Errorf("database error when checking the existence in executed order history")
	}
	exists_in_cancel_history, err = cancelledOrderExists(conn, orderId)
	if err != nil {
		return []OpenOrderTuple{}, []ExecutedOrderHistoryTuple{}, []CancelledOrderHistoryTuple{}, fmt.Errorf("database error when checking the existence in cancelled order history")
	}
	exists_in_open_orders, err = checkOrderExists(conn, orderId)
	if err != nil {
		return []OpenOrderTuple{}, []ExecutedOrderHistoryTuple{}, []CancelledOrderHistoryTuple{}, fmt.Errorf("database error when checking the existence in open order history")
	}

	if !exists_in_open_orders && !exists_in_executed_history && !exists_in_cancel_history {
		return []OpenOrderTuple{}, []ExecutedOrderHistoryTuple{}, []CancelledOrderHistoryTuple{}, fmt.Errorf("no such order exists")
	}

	var executedOrderHistoryQueryResult []ExecutedOrderHistoryTuple
	if exists_in_executed_history {
		var executed_history_node_list []string
		executed_history_node_list, err = GetExecutedOrderSliceList(conn, orderId)
		if err != nil {
			return []OpenOrderTuple{}, []ExecutedOrderHistoryTuple{}, []CancelledOrderHistoryTuple{}, fmt.Errorf("database error when retrieving the executed order history")
		}
		executedOrderHistoryQueryResult = parseExcutedHistoryNodeList(executed_history_node_list)
	}

	var openOrderQueryResult []OpenOrderTuple
	if exists_in_open_orders {
		var amount float64
		amount, err = GetOrderAmount(conn, orderId)
		if err != nil {
			return []OpenOrderTuple{}, []ExecutedOrderHistoryTuple{}, []CancelledOrderHistoryTuple{}, fmt.Errorf("database error when retrieving the open order")
		}
		var symbolName_n_orderType []string
		symbolName_n_orderType, err = GetSymbolNameAndOrderType(conn, orderId)
		if err != nil {
			return []OpenOrderTuple{}, []ExecutedOrderHistoryTuple{}, []CancelledOrderHistoryTuple{}, fmt.Errorf("database error when retrieving the open order")
		}
		if symbolName_n_orderType[1] == "sell" {
			amount = -amount
		}
		amount_in_string := fmt.Sprintf("%f", amount)
		openOrderQueryResult = append(openOrderQueryResult, OpenOrderTuple{CurrentAmount: amount_in_string})
	}

	var cancelledOrderHistoryQueryResult []CancelledOrderHistoryTuple

	if exists_in_cancel_history {
		var amount float64
		var time string
		amount, time, err = getAmountAndTimeForCancelledOrderFromCancelHistory(conn, orderId)
		if err != nil {
			return []OpenOrderTuple{}, []ExecutedOrderHistoryTuple{}, []CancelledOrderHistoryTuple{}, fmt.Errorf("database error when retrieving the cancelled order history")
		}
		amount_in_string := fmt.Sprintf("%f", amount)
		cancelledOrderHistoryQueryResult = append(cancelledOrderHistoryQueryResult, CancelledOrderHistoryTuple{CancelledAmount: amount_in_string, CancelledTime: time})
	}

	return openOrderQueryResult, executedOrderHistoryQueryResult, cancelledOrderHistoryQueryResult, nil
}

/*
		MatchOrder will match open order with orderId with possible open orders.
		If matched, a transaction is executed automatically, and an executed history is inserted.
		This function will match as many time as possible.
		The order will be removed if it becomes an empty(order amount = 0) order after transactions.
		Orders which are matched by this function will be removed if they become empty orders after a transaction.
		The removals are done in executeMatch function, which is called after finding a match.
		matchForBuyOrder and matchForSellOrder are sub functions to implement MatchOrder's functionality.
		Their inputs are the same as MatchOrder, and their logic is described as above.
	input --
		orderId: represents the order that you want to find match for
		uid: account id associated with the order with order id
		symbolName: symbol name of the order with order id
		limitPrice: limit price of the order with order id
		amount: order amount of the order with order id
		orderType: order type(buy/sell) of the order with order id
	output --
		this function will NOT check if the order is an open order and has been added to openOrderBook,
		so MAKE SURE you are matching an open order.

		err:
		database err
*/
func MatchOrder(conn *redigo.Conn, orderId string, uid string, symbolName string, limitPrice float64, amount float64, orderType string) error {
	if orderType == ORDER_TYPE_BUY {
		err := matchForBuyOrder(conn, orderId, uid, symbolName, limitPrice, amount)
		if err != nil {
			return err
		}
	} else {
		err := matchForSellOrder(conn, orderId, uid, symbolName, limitPrice, amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func matchForBuyOrder(conn *redigo.Conn, orderId string, uid string, symbolName string, limitPrice float64, amount float64) error {
	buyOrderId := orderId
	for {
		empty, err := isSellOrderBookEmpty(conn, symbolName)
		if err != nil {
			return fmt.Errorf("database error when checking a sell order book is empty")
		}

		if empty {
			return nil
		}

		var sell_order_id_with_min_price string
		var sell_order_min_price float64
		sell_order_id_with_min_price, sell_order_min_price, err = peekSellOrderWithMinPriceInSellOrdrerBook(conn, symbolName)
		if err != nil {
			return fmt.Errorf("database error when peeking the min price in a sell order book")
		}

		buy_order_price := limitPrice
		if sell_order_min_price > buy_order_price {
			return nil
		}

		err = executeMatch(conn, buyOrderId, sell_order_id_with_min_price, symbolName, "buy")
		if err != nil {
			return err
		}

		var exists bool
		exists, err = checkOrderExists(conn, orderId)

		if !exists {
			return nil
		}
	}
}

func matchForSellOrder(conn *redigo.Conn, orderId string, uid string, symbolName string, limitPrice float64, amount float64) error {
	sellOrderId := orderId
	for {
		empty, err := isBuyOrderBookEmpty(conn, symbolName)
		if err != nil {
			return fmt.Errorf("database error when checking a buy order book is empty")
		}

		if empty {
			return nil
		}

		var buy_order_id_with_max_price string
		var buy_order_max_price float64
		buy_order_id_with_max_price, buy_order_max_price, err = peekBuyOrderWithMaxPriceInBuyOrdrerBook(conn, symbolName)
		if err != nil {
			return fmt.Errorf("database error when peeking the max price in a buy order book")
		}

		sell_order_price := limitPrice
		if sell_order_price > buy_order_max_price {
			return nil
		}

		err = executeMatch(conn, buy_order_id_with_max_price, sellOrderId, symbolName, "sell")
		if err != nil {
			return err
		}

		var exists bool
		exists, err = checkOrderExists(conn, orderId)

		if !exists {
			return nil
		}
	}
}

/*
		executeMatch does a transaction between a matched buy/sell orders pair.
		An executed history is added into executed history when execution is done.
		This function will NOT check if the orders match or not. The user has to make sure two orders match.
		Since executed history does not contain info about orderType(buy/sell), we set transaction amount in executed hitory to negative as "sell"
		This function will NOT validate both orders existence and openness
		This function will atomatically remove orders when an order's amount become 0(empty order).
		This function will also remove the order which inits the transaction when it becomes empty.
	input --
		buyOrderId: buy order's id
		sellOrderId: sell order's id
		symbolName: the symbol name of buy and sell order
		transInitOrderType: the order type of the order which inits this execution
		eg:
		if a buy order is added to the market and the matching engine matches it with a existed sell order,
		the buy order is the order to init this execution
	output --
		err:
		database error

*/
func executeMatch(conn *redigo.Conn, buyOrderId string, sellOrderId string, symbolName string, transInitOrderType string) error {
	var sell_order_amount, buy_order_amount float64
	var err error
	sell_order_amount, err = GetOrderAmount(conn, sellOrderId)
	if err != nil {
		return fmt.Errorf("database error when retrieving the sell order's amount")
	}
	buy_order_amount, err = GetOrderAmount(conn, buyOrderId)
	if err != nil {
		return fmt.Errorf("database error when retrieving the buy order's amount")
	}

	transaction_amount := math.Min(buy_order_amount, sell_order_amount)

	var buyer_uid, seller_uid string
	buyer_uid, err = GetOrderUid(conn, buyOrderId)
	if err != nil {
		return fmt.Errorf("database error when retrieving the buy order's uid")
	}
	seller_uid, err = GetOrderUid(conn, sellOrderId)
	if err != nil {
		return fmt.Errorf("database error when retrieving the sell order's uid")
	}

	var buy_order_limit_price, sell_order_limit_price float64
	sell_order_limit_price, err = GetOrderLimitPrice(conn, sellOrderId)
	if err != nil {
		return fmt.Errorf("database error when retrieving the sell order's limit price")
	}
	buy_order_limit_price, err = GetOrderLimitPrice(conn, buyOrderId)
	if err != nil {
		return fmt.Errorf("database error when retrieving the buy order's limit price")
	}

	var transaction_price float64
	if transInitOrderType == "buy" {
		transaction_price = sell_order_limit_price
	} else {
		transaction_price = buy_order_limit_price
	}

	_, err = increaseSymbolPosition(conn, buyer_uid, symbolName, transaction_amount)
	if err != nil {
		return fmt.Errorf("database error when adding symbol to the buyer's account")
	}

	_, err = increaseAccountBalance(conn, seller_uid, transaction_price*transaction_amount)
	if err != nil {
		return fmt.Errorf("database error when adding balance to the seller's account")
	}

	refundToBuyer := (buy_order_limit_price - transaction_price) * transaction_amount
	if refundToBuyer > 0 {
		_, err = increaseAccountBalance(conn, buyer_uid, refundToBuyer)
		if err != nil {
			return fmt.Errorf("database error when refunding balance to the buyer's account")
		}
	}

	if transaction_amount == buy_order_amount {
		err = removeBuyOrderFromBuyOrderBook(conn, symbolName, buyOrderId)
		if err != nil {
			return fmt.Errorf("database error when removing empty order from buy order book")
		}
		err = removeOrder(conn, buyOrderId)
		if err != nil {
			return fmt.Errorf("database error when removing empty buy order from buy orders")
		}

	} else {
		_, err = decreaseOrderAmount(conn, buyOrderId, transaction_amount)
		if err != nil {
			return fmt.Errorf("database error when decrease amount from buy order")
		}
	}

	if transaction_amount == sell_order_amount {
		err = removeSellOrderFromSellOrderBook(conn, symbolName, sellOrderId)
		if err != nil {
			return fmt.Errorf("database error when removing empty order from sell order book")
		}
		err = removeOrder(conn, sellOrderId)
		if err != nil {
			return fmt.Errorf("database error when removing empty buy order from sell orders")
		}
	} else {
		_, err = decreaseOrderAmount(conn, sellOrderId, transaction_amount)
		if err != nil {
			return fmt.Errorf("database error when decrease amount from sell order")
		}
	}

	current_time := getCurrentTimeInString()
	err = InsertExcutedOrderToExcutedHistory(conn, buyOrderId, transaction_amount, transaction_price, current_time)
	if err != nil {
		return fmt.Errorf("database error when inserting buy order executed history")
	}
	// Since executed history does not contain info about orderType(buy/sell), we set transaction amount in executed hitory to negative as "sell"
	err = InsertExcutedOrderToExcutedHistory(conn, sellOrderId, -transaction_amount, transaction_price, current_time)
	if err != nil {
		return fmt.Errorf("database error when inserting sell order executed history")
	}

	return nil
}
