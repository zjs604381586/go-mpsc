# Go-MPSC

A **high-performance MPSC queue** between [wait-free](https://en.wikipedia.org/wiki/Non-blocking_algorithm) and [lock-free](https://en.wikipedia.org/wiki/Non-blocking_algorithm), that is a multi-producer but one consumer queue

# Time-consuming test with using go channel:

**Test Case:** 

  1W producers, each producer writes 100 times. Statistics on the time taken by all producers to write to the queue
  
**Go version**: 
  
  go version go1.16.11 darwin/amd64.    
   
**Machine information**：   
  * Model: macOS Big Sur 11.5.2       
  * Cpu: Six cores Intel Core i7   
  * Memory: 16g
        
**Test Result:**
  

   
  ![image](https://user-images.githubusercontent.com/17305630/159618064-3e4fcd10-3440-494b-bc07-54a5777fe73a.png)
