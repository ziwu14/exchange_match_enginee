# Danger Log

1. Atomicity

   Some operations should be considered as a transaction, so that no operations will be executed if one of them fails. One example is when two orders matches, the engine should add balance to the seller and add symbols to the buyer.

2. Critical sections

   ```
   # read from db
   Critical
   	# check the value from db
   	# branch1 : do stuff in db
   	# branch2 : do other stuff in db
   Critical
   ...
   ```

3. Lock granularity

    For simplicity, we are using a global read-write lock in our matching engine. We are wondering whether a better concurrency pattern can be used in this case.