#!/bin/bash
./buy_test1.sh
./sell_test1.sh
./buy_test2.sh
./sell_test2.sh
./sell_test3.sh
./buy_test3.sh
cat query1.txt | nc localhost 12345
cat cancel1.txt | nc localhost 12345
cat query1.txt | nc localhost 12345
