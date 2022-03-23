# MPSC-Go

A **high-performance MPSC queue** between [wait-free](https://en.wikipedia.org/wiki/Non-blocking_algorithm) and [lock-free](https://en.wikipedia.org/wiki/Non-blocking_algorithm), that is, a multi-producer but consumer queue

# Time-consuming test with using channel:

Test Case: 

    1W producers, each producer writes 100 times. Statistics on the time taken by all producers to write to the queue
  
Test Result:
  
