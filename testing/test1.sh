#!/bin/bash
cat create1.txt | nc localhost 12345 # create buyer, id=12345 $10000
cat create2.txt | nc localhost 12345 # create seller, id=34567 bitcoin 100
cat buy2.txt | nc localhost 12345 # buyer buy 50 bitcoin at $7
cat sell2.txt | nc localhost 12345 # seller sell 10 bitcoin at $4
cat sell3.txt | nc localhost 12345 # seller sell 10 bitcoin at $5
cat query_before_cancel.txt | nc localhost 12345 # query buy order
cat cancel.txt | nc localhost 12345 # cancel buy order
cat query_before_cancel.txt | nc localhost 12345 # query buy order
