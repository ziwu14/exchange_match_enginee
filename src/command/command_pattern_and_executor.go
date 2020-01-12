package command

import (
	"app/businessLogic"
	"fmt"
	"strconv"
	"sync"
	"app/uniqueKeyGenerator"

	redigo "github.com/gomodule/redigo/redis"
)

type Command interface {
	execute(pool *redigo.Pool, readWriteLock *sync.RWMutex)
	getResponse() string
}

type CreateAccoutCommand struct {
	Uid     string
	Balance float64

	Err      error
	Response string
}

func (c *CreateAccoutCommand) execute(pool *redigo.Pool, readWriteLock *sync.RWMutex) {
	readWriteLock.Lock()
	defer readWriteLock.Unlock()

	Err := businessLogic.CreateAccount(pool, c.Uid, c.Balance)
	if Err != nil {
		c.Response = fmt.Sprintf("<error id=\"%s\">%s</error>", c.Uid, Err)
		return
	} else {
		c.Response = fmt.Sprintf("<created id=\"%s\"/>", c.Uid)
	}
}

func (c *CreateAccoutCommand) getResponse() string {
	return c.Response
}

type SetOrAddSymbolPositionToAccountCommand struct {
	Uid        string
	SymbolName string
	Amount     float64

	Err      error
	Response string
}

func (c *SetOrAddSymbolPositionToAccountCommand) execute(pool *redigo.Pool, readWriteLock *sync.RWMutex) {
	readWriteLock.Lock()
	defer readWriteLock.Unlock()

	Err := businessLogic.SetOrAddSymbolPositionToAccount(pool, c.Uid, c.SymbolName, c.Amount)
	if Err != nil {
		c.Response = fmt.Sprintf("<error sym=\"%s\" id=\"%s\">%s</error>", c.SymbolName, c.Uid, Err)
		return
	} else {
		c.Response = fmt.Sprintf("<created sym=\"%s\" id=\"%s\"/>", c.SymbolName, c.Uid)
	}
}

func (c *SetOrAddSymbolPositionToAccountCommand) getResponse() string {
	return c.Response
}

type SetBuyOrderCommand struct {
	OrderId    string
	Uid        string
	SymbolName string
	LimitPrice float64
	Amount     float64

	Err      error
	Response string
}

func (c *SetBuyOrderCommand) execute(pool *redigo.Pool, readWriteLock *sync.RWMutex) {
	readWriteLock.Lock()
	defer readWriteLock.Unlock()

	orderId, err := uniqueKeyGenerator.GetNewOrderId(pool)
	c.OrderId = strconv.Itoa(orderId)

	if err != nil {
		c.Response = fmt.Sprintf("<error sym=\"%s\" Amount=\"%.2f\" limit=\"%.2f\" >%s</error>", c.SymbolName, c.Amount, c.LimitPrice, "error when generating orderId")
		return
	}

	err = businessLogic.SetBuyOrder(pool, c.OrderId, c.Uid, c.SymbolName, c.LimitPrice, c.Amount)
	if err != nil {
		c.Response = fmt.Sprintf("<error sym=\"%s\" Amount=\"%.2f\" limit=\"%.2f\" >%s</error>", c.SymbolName, c.Amount, c.LimitPrice, err)
		return
	} else {
		c.Response = fmt.Sprintf("<opened sym=\"%s\" Amount=\"%.2f\" limit=\"%.2f\" id=\"%s\"/>", c.SymbolName, c.Amount, c.LimitPrice, c.OrderId)
	}
}

func (c *SetBuyOrderCommand) getResponse() string {
	return c.Response
}

type SetSellOrderCommand struct {
	OrderId    string
	Uid        string
	SymbolName string
	LimitPrice float64
	Amount     float64

	Err      error
	Response string
}

func (c *SetSellOrderCommand) execute(pool *redigo.Pool, readWriteLock *sync.RWMutex) {
	readWriteLock.Lock()
	defer readWriteLock.Unlock()

	orderId, err := uniqueKeyGenerator.GetNewOrderId(pool)
	c.OrderId = strconv.Itoa(orderId)

	if err != nil {
		c.Response = fmt.Sprintf("<error sym=\"%s\" Amount=\"%.2f\" limit=\"%.2f\" >%s</error>", c.SymbolName, -c.Amount, c.LimitPrice, "error when generating orderId")
		return
	}

	err = businessLogic.SetSellOrder(pool, c.OrderId, c.Uid, c.SymbolName, c.LimitPrice, c.Amount)
	if err != nil {
		c.Response = fmt.Sprintf("<error sym=\"%s\" Amount=\"%.2f\" limit=\"%.2f\" >%s</error>", c.SymbolName, -c.Amount, c.LimitPrice, err)
		return
	} else {

		c.Response = fmt.Sprintf("<opened sym=\"%s\" Amount=\"%.2f\" limit=\"%.2f\" id=\"%s\"/>", c.SymbolName, -c.Amount, c.LimitPrice, c.OrderId)
	}
}

func (c *SetSellOrderCommand) getResponse() string {
	return c.Response
}

type CancelOpenOrderCommand struct {
	OrderId string

	Err      error
	Response string
}

func (c *CancelOpenOrderCommand) execute(pool *redigo.Pool, readWriteLock *sync.RWMutex) {
	readWriteLock.Lock()
	defer readWriteLock.Unlock()

	Err_in_cancel := businessLogic.CancelOpenOrder(pool, c.OrderId)

	if Err_in_cancel != nil {
		c.Response = fmt.Sprintf("<error id=\"%s\">%s</error>", c.OrderId, Err_in_cancel)
		return
	}
	_, executedOrderHistory, cancelledOrderHistory, Err_in_query := businessLogic.QueryOrderStatusAndHistory(pool, c.OrderId)
	if Err_in_query != nil {
		c.Response = fmt.Sprintf("<error id=\"%s\">%s</error>", c.OrderId, Err_in_query)
		return
	} else {
		var executedHistoryResponse string
		for _, executedHistoryTuple := range executedOrderHistory {
			executedHistoryResponse +=
				fmt.Sprintf("  <executed shares=%s price=%s time=%s/>",
					executedHistoryTuple.TransactionAmount,
					executedHistoryTuple.TransactionPrice,
					executedHistoryTuple.TransactionTime) + "\n"
		}

		c.Response =
			fmt.Sprintf("<canceled id=\"%s\">", c.OrderId) + "\n" +
				fmt.Sprintf("  <canceled shares=-%s time=%s/>",
					cancelledOrderHistory[0].CancelledAmount,
					cancelledOrderHistory[0].CancelledTime) + "\n" +
				executedHistoryResponse +
				fmt.Sprintf("</canceled>")
	}
}

func (c *CancelOpenOrderCommand) getResponse() string {
	return c.Response
}

type QueryOrderStatusAndHistoryCommand struct {
	OrderId string

	Err      error
	Response string
}

func (c *QueryOrderStatusAndHistoryCommand) execute(pool *redigo.Pool, readWriteLock *sync.RWMutex) {
	readWriteLock.RLock()
	defer readWriteLock.RUnlock()

	openOrderTuples, executedOrderHistory, cancelledOrderHistory, Err_in_query := businessLogic.QueryOrderStatusAndHistory(pool, c.OrderId)
	if Err_in_query != nil {
		c.Response = fmt.Sprintf("<error id=\"%s\">%s</error>", c.OrderId, Err_in_query)
		return
	} else {
		var cancelledHistoryResponse string
		if len(cancelledOrderHistory) > 0 {
			cancelledHistoryResponse =
				fmt.Sprintf("  <canceled shares=%s time=%s/>",
					cancelledOrderHistory[0].CancelledAmount,
					cancelledOrderHistory[0].CancelledTime) + "\n"
		}

		var executedHistoryResponse string
		if len(executedOrderHistory) > 0 {
			for _, executedHistoryTuple := range executedOrderHistory {
				executedHistoryResponse +=
					fmt.Sprintf("  <executed shares=%s price=%s time=%s/>",
						executedHistoryTuple.TransactionAmount,
						executedHistoryTuple.TransactionPrice,
						executedHistoryTuple.TransactionTime) + "\n"
			}
		}

		var openOrderTupleResponse string
		if len(openOrderTuples) > 0 {
			openOrderTupleResponse =
				fmt.Sprintf("  <opened shares=%s/>",
					openOrderTuples[0].CurrentAmount) + "\n"
		}

		c.Response =
			fmt.Sprintf("<status id=\"%s\">", c.OrderId) + "\n" +
				openOrderTupleResponse + cancelledHistoryResponse + executedHistoryResponse +
				fmt.Sprintf("</status>")
	}
}

func (c *QueryOrderStatusAndHistoryCommand) getResponse() string {
	return c.Response
}

type CommandListExecutor struct {
	Pool     *redigo.Pool
	Response string

	ReadWriteLock *sync.RWMutex
}

func (e *CommandListExecutor) Execute(commandList []Command) {
	e.Response += "<results>" + "\n"
	for _, command := range commandList {
		command.execute(e.Pool, e.ReadWriteLock)
		e.Response += "  " + command.getResponse() + "\n"
	}
	e.Response += "</results>" + "\n"
}

func (e *CommandListExecutor) GetResponse() string {
	return e.Response
}

// func main() {
// 	pool := redis.NewRConnectionPool(
// 		redis.Config{
// 			Server:              "redis:6379",
// 			Password:            "",
// 			MaxIdle:             100,
// 			MaxActive:           12000,
// 			IdleTimeout:         240 * time.Second,
// 			KEY_PREFIX:          "",
// 			KEY_DELIMITER:       "",
// 			KEY_VAR_PLACEHOLDER: "",
// 		},
// 	)

// 	connection := pool.Get()
// 	defer connection.Close()
// 	conn := (&connection)

// 	redis.FlushAll(conn)

// 	c1 := &CreateAccoutCommand{Uid: "12345", Balance: 10000}
// 	c1.execute(pool)
// 	fmt.Println(c1.getResponse())
// 	// c1.execute(pool)
// 	// fmt.Println(c1.getResponse())

// 	c2 := &SetOrAddSymbolPositionToAccountCommand{Uid: "12345", SymbolName: "bitcoin", Amount: 100}
// 	c2.execute(pool)
// 	fmt.Println(c2.getResponse())
// 	// c2.execute(pool)
// 	// fmt.Println(c2.getResponse())
// 	// c2.Uid = "non-exist"
// 	// c2.execute(pool)
// 	// fmt.Println(c2.getResponse())

// 	d1 := &CreateAccoutCommand{Uid: "111", Balance: 0}
// 	d1.execute(pool)
// 	d2 := &SetOrAddSymbolPositionToAccountCommand{Uid: "111", SymbolName: "bitcoin", Amount: 100}
// 	d2.execute(pool)
// 	d3 := &SetSellOrderCommand{OrderId: "101", Uid: "111", SymbolName: "bitcoin", LimitPrice: 100, Amount: 10}
// 	d3.execute(pool)

// 	q2 := &QueryOrderStatusAndHistoryCommand{OrderId: "101"}
// 	q2.execute(pool)
// 	fmt.Println(q2.getResponse())

// 	c3 := &SetBuyOrderCommand{OrderId: "1", Uid: "12345", SymbolName: "bitcoin", LimitPrice: 100, Amount: 100}
// 	c3.execute(pool)
// 	fmt.Println(c3.getResponse())

// 	q1 := &QueryOrderStatusAndHistoryCommand{OrderId: "1"}
// 	q1.execute(pool)
// 	fmt.Println(q1.getResponse())

// 	c5 := &CancelOpenOrderCommand{OrderId: "1"}
// 	c5.execute(pool)
// 	fmt.Println(c5.getResponse())

// 	q1.execute(pool)
// 	fmt.Println(q1.getResponse())

// 	commandList := []Command{c1, c2, d1, d2, d3, q2, c3, q1, c5, q1}
// 	commandListExecutor := &CommandListExecutor{Response: "", Pool: pool}
// 	commandListExecutor.Execute(commandList)
// 	fmt.Println("executor result:\n")
// 	fmt.Println(commandListExecutor.GetResponse())

// 	// c4 := &SetSellOrderCommand{OrderId: "1", Uid: "12345", SymbolName: "bitcoin", LimitPrice: 100, Amount: 100}
// 	// c4.execute(pool)
// 	// fmt.Println(c4.getResponse())

// 	// c4.Amount = 10000
// 	// c4.execute(pool)
// 	// fmt.Println(c4.getResponse())

// 	// c6 := &QueryOrderStatusAndHistoryCommand{}

// 	// commandList := []Command{c1, c2, c3, c4, c5, c6}
// 	// for _, command := range commandList {
// 	// 	command.execute()
// 	// }

// 	// var commandResponseList []string
// 	// for _, command := range commandList {
// 	// 	commandResponseList = append(commandResponseList, command.getResponse())
// 	// }

// 	// fmt.Println(commandResponseList)
// }
