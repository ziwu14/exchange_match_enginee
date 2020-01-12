time seq 2000 | parallel -n0 "cat query1.txt | nc localhost 12345" > /dev/null
