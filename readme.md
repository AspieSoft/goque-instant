# Goque Instant

[![donation link](https://img.shields.io/badge/buy%20me%20a%20coffee-paypal-blue)](https://paypal.me/shaynejrtaylor?country.x=US&locale.x=en_US)

A fast and concurrent object queue for golang.

This module avoids shifting indexes and memory in an array.
Instead it utilizes how a uint16 works, where it goes back to 0 when it passes its maximum size.
This will reset the number to loop back to the front of the queue.
When this happens, the queue is checked for an empty slot, and if unavailable will append to an overflow queue.
The overflow queue will concurrently be added to the main queue when a new spot is available.

Notice: This is not for objects where the order is important. It should maintain a consistant order, but objects are added to the queue concurrently, which could result in an offset queue order.

## Installation

```shell script
go get github.com/AspieSoft/goque-instant
```

## Usage

```go

import (
  "github.com/AspieSoft/goque-instant"
)

func main(){
  myQueue := goque.New[int /* any */]()

  myQueue.Add(1)
  myQueue.Add(2)
  myQueue.Add(3)

  fmt.Println(myQueue.Next()) // 1
  fmt.Println(myQueue.Next()) // 2
  fmt.Println(myQueue.Next()) // 3

  // waits for the queue to reach 0
  myQueue.Wait()

  if myQueue.Len() == 0 {
    myQueue.Stop()
  }

  myq.WaitAndStop()
}

```
