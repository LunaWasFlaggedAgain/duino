# duino
A simple duinocoin pool API wrapper for golang.

# Benchmarks
Benchmarked using [this code](https://gist.github.com/LunaWasFlaggedAgain/14f54734db38c50e634532c30ddf2419)
## AMD Ryzen 5 5600G (12) @ 3.900GHz
```
DoJob(AVR): 74 (81.33µs)
DoJobMulti(AVR): 74 (812.57µs)
DoJob(ESP32): 52020 (18.509911ms)
DoJobMulti(ESP32): 52020 (5.42038ms)
DoJob(LOW): 2520154 (1.620081394s)
DoJobMulti(LOW): 2520154 (284.422773ms)
DoJob(MEDIUM): 14103739 (8.843125167s)
DoJobMulti(MEDIUM): 14103739 (1.384273916s)
DoJob(HIGH): 59515575 (42.686280546s)
DoJobMulti(HIGH): 59515575 (6.353000107s)
```