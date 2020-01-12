# README

Description: 

This is an exchange matching engine written in Go. We use Redis as the database.

Usage:

1. Git clone the repository to a workspace.

2. On linux, set env `COMPOSE_HTTP_TIMEOUT=3600`, otherwise docker container might exit due to a http request time out error

3. Run `sudo docker-compose build && sudo docker-compose up`, port 12345 will be used by the matching engine. 

4. Inside folder *testing*:

   *test1.sh* is a sequential test to verify the correctness of the matching engine.

   *test2.sh* is a parallel test to verify the correctness.  

   You have to install *parallel* to run test2.sh: `sudo apt-get install parallel` (on ubuntu)

5. *test1.sh*'s testcase: 

   ```
   after run test1.sh
   enter redis container and run redis-cli, then run
   HGET account:12345 balance       (get the balance of buyer:12345)
   HGET account:34567 balance       (get the balance of seller:34567)
   HGET account:12345:SPY amount    (buyer's SPY amount)
   HGET account:34567:SPY amount   (seller's SPY amount)
   
   final state should be:
   - buyer balance = $9860, SPY = 20
   - seller balance = $140, SPY = 80
   
   ```

   

   * buyer-uid: 12345, balance = $10000
   * seller-uid: 34567, SPY = 100
   * set buy order: 
     * orderid = 1, amount = 50, limit = $7
   * set sell orders:
     * orderid = 2, amount = -10, limit = $4 (matches with 1)
     * orderid = 3, amount = -10, limit = $5 (matches with 1)
   * open order: orderid = 1, amount = 30, limit = $7 ( with 2 executed history, buy 10 SPY at 7 dollars )
   * cancel open order: orderid = 1, refund $210 to buyer

6. *test2.sh*'s testcase:

   ```
   after run test2.sh
   enter redis container and run redis-cli, then run
   HGET account:12345 balance       (get the balance of buyer:12345)
   HGET account:34567 balance       (get the balance of seller:34567)
   HGET account:12345:SPY amount    (buyer's SPY amount)
   HGET account:34567:SPY amount   (seller's SPY amount)
   
   final state should be:
   - buyer balance = $0, SPY = 100
   - seller balance = $10000, SPY = 0
   ```

   

   * create buyer-uid: 12345, balance = $10000 for 100 times 
   * create seller-uid: 34567, with SPY = 1 each time (add 100 SPY in total) for 100 times
   * create buy order 200 times : buy 1 SPY, $100/ each (100 success, 100 insufficient fund)
   * create sell order 200 times: sell 1 SPY, $100/ each (100 success, 100 insufficient symbol) (all 100 sell orders match successfully)