package xmlParser

import (
	cmd "app/command"
	"fmt"
	"strconv"

	"github.com/beevik/etree"
)

type Parser interface {
	Parse(string) ([]cmd.Command, error)
}

type XmlParser struct{}

func (parser *XmlParser) Parse(xml string) ([]cmd.Command, error) {
	request := etree.NewDocument()

	if err := request.ReadFromString(xml); err != nil {
		return []cmd.Command{}, fmt.Errorf("xml format error from etree")
	}

	var commandList []cmd.Command

	createElement := request.SelectElement("create")
	if createElement != nil {
		// create
		for _, req := range createElement.ChildElements() {
			if req.Tag == "account" {
				uid, balance_in_string := readElementWith2Attr(req, "id", "balance")
				if uid == "" || balance_in_string == "" {
					return []cmd.Command{}, fmt.Errorf("xml format error")
				}

				balance, err := strconv.ParseFloat(balance_in_string, 64)
				if err != nil {
					return []cmd.Command{}, fmt.Errorf("xml format error")
				}
				commandList = append(commandList,
					&cmd.CreateAccoutCommand{
						Uid:     uid,
						Balance: balance})
			} else if req.Tag == "symbol" {
				symbolName := readElementWith1Attr(req, "sym")
				if symbolName == "" {
					return []cmd.Command{}, fmt.Errorf("xml format error")
				}
				if len(req.SelectElements("account")) != len(req.ChildElements()) {
					return []cmd.Command{}, fmt.Errorf("xml format error")
				}
				// parse account inside symbol
				for _, innerReq := range req.ChildElements() {
					uid := readElementWith1Attr(innerReq, "id")
					if uid == "" {
						return []cmd.Command{}, fmt.Errorf("xml format error")
					}
					amount_in_string := innerReq.Text()
					amount, err := strconv.ParseFloat(amount_in_string, 64)
					if err != nil {
						return []cmd.Command{}, fmt.Errorf("xml format error")
					}
					commandList = append(commandList,
						&cmd.SetOrAddSymbolPositionToAccountCommand{
							Uid:        uid,
							SymbolName: symbolName,
							Amount:     amount})
				}
			}
		}
	}

	transactionElement := request.SelectElement("transactions")
	if transactionElement != nil {
		// transaction
		if len(transactionElement.ChildElements()) == 0 {
			return []cmd.Command{}, fmt.Errorf("xml format error")
		}

		uid := readElementWith1Attr(transactionElement, "id")
		if uid == "" {
			return []cmd.Command{}, fmt.Errorf("xml format error")
		}

		for _, req := range transactionElement.ChildElements() {
			if req.Tag == "order" {
				symbolName, amount_in_string, limitPrice_in_string := readElementWith3Attr(req, "sym", "amount", "limit")
				if symbolName == "" || amount_in_string == "" || limitPrice_in_string == "" {
					return []cmd.Command{}, fmt.Errorf("xml format error")
				}

				var amount, limitPrice float64
				var err error
				amount, err = strconv.ParseFloat(amount_in_string, 64)
				if err != nil {
					return []cmd.Command{}, fmt.Errorf("xml format error")
				}

				limitPrice, err = strconv.ParseFloat(limitPrice_in_string, 64)
				if err != nil {
					return []cmd.Command{}, fmt.Errorf("xml format error")
				}

				//TODO: Moved orderId generating to specified commands
				if amount < 0 {
					commandList = append(commandList,
						&cmd.SetSellOrderCommand{
							Uid:        uid,
							SymbolName: symbolName,
							LimitPrice: limitPrice,
							Amount:     -amount})
				}

				if amount > 0 {
					commandList = append(commandList,
						&cmd.SetBuyOrderCommand{
							Uid:        uid,
							SymbolName: symbolName,
							LimitPrice: limitPrice,
							Amount:     amount})
				}
			} else if req.Tag == "query" {
				orderId := readElementWith1Attr(req, "id")
				if orderId == "" {
					return []cmd.Command{}, fmt.Errorf("xml format error")
				}

				commandList = append(commandList,
					&cmd.QueryOrderStatusAndHistoryCommand{
						OrderId: orderId})
			} else if req.Tag == "cancel" {
				orderId := readElementWith1Attr(req, "id")
				if orderId == "" {
					return []cmd.Command{}, fmt.Errorf("xml format error")
				}

				commandList = append(commandList,
					&cmd.CancelOpenOrderCommand{
						OrderId: orderId})
			} else {
				return []cmd.Command{}, fmt.Errorf("xml format error")
			}
		}
	}

	if createElement == nil && transactionElement == nil {
		return []cmd.Command{}, fmt.Errorf("xml format error")
	}

	return commandList, nil
}

func attrExists(element *etree.Element, key string) bool {
	value := element.SelectAttrValue(key, "")
	if value == "" {
		return false
	}
	return true
}

func readElementWith2Attr(element *etree.Element, key1 string, key2 string) (value1 string, value2 string) {
	if !attrExists(element, key1) || !attrExists(element, key2) {
		return "", ""
	}
	value1 = element.SelectAttrValue(key1, "")
	value2 = element.SelectAttrValue(key2, "")
	return
}

func readElementWith3Attr(element *etree.Element, key1 string, key2 string, key3 string) (value1 string, value2 string, value3 string) {
	if !attrExists(element, key1) || !attrExists(element, key2) || !attrExists(element, key3) {
		return "", "", ""
	}
	value1 = element.SelectAttrValue(key1, "")
	value2 = element.SelectAttrValue(key2, "")
	value3 = element.SelectAttrValue(key3, "")
	return
}

func readElementWith1Attr(element *etree.Element, key string) (value string) {
	if !attrExists(element, key) {
		return ""
	}
	value = element.SelectAttrValue(key, "")
	return
}

// func main() {
// 	b, _ := ioutil.ReadFile("xml_example_5.xml")
// 	xml := string(b)

// 	var parser XmlParser
// 	commandList, err := parser.parse(xml)
// 	if err != nil {
// 		fmt.Println("<err>:", err)
// 	} else {
// 		for _, command := range commandList {
// 			command.execute()
// 		}
// 	}
// }
