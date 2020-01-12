package businessLogic

import (
	"strconv"
	"time"
)

func isBase10NumberSequense(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func getCurrentTimeInString() string {
	currentTime := time.Now()
	epoch := currentTime.Unix()
	epochInString := strconv.FormatInt(epoch, 10)
	return epochInString
}

func parseExcutedHistoryNodeList(executedHistoryNodeList []string) []ExecutedOrderHistoryTuple {
	numberOfTuples := len(executedHistoryNodeList) / 3
	var tupleList []ExecutedOrderHistoryTuple
	for i := 0; i < numberOfTuples; i++ {
		tuple := ExecutedOrderHistoryTuple{
			TransactionAmount: executedHistoryNodeList[i*3],
			TransactionPrice:  executedHistoryNodeList[i*3+1],
			TransactionTime:   executedHistoryNodeList[i*3+2],
		}
		tupleList = append(tupleList, tuple)
	}
	return tupleList
}
