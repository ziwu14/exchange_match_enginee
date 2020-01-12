#!/bin/bash
time seq 10000 | parallel -n0 "cat create1.txt | nc localhost 12345"  > results/c1p_r.txt
time seq 10000 | parallel -n0 "cat create3p.txt | nc localhost 12345"  > results/c2p_r.txt
time seq 20000 | parallel -n0 "cat buy5p.txt | nc localhost 12345"  > results/b5p_r.txt
time seq 20000 | parallel -n0 "cat sell6p.txt | nc localhost 12345"  > results/s6p_r.txt
#seq 200 | parallel -n0 "cat buy6p.txt | nc localhost 12345" > results/b6p_r.txt
#seq 200 | parallel -n0 "cat sell7p.txt | nc localhost 12345" > results/s7p_r.txt
