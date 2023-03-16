# Goque Instant

[![donation link](https://img.shields.io/badge/buy%20me%20a%20coffee-paypal-blue)](https://paypal.me/shaynejrtaylor?country.x=US&locale.x=en_US)

A fast and concurrent object queue for golang.

This module avoids shifting indexes and memory in an array.
Instead it utilizes how a uint16 works, where it goes back to 0 when it passes its maximum size.
This will reset the number to loop back to the front of the queue.
When this happens, the queue is checked for an empty slot, and if unavailable will append to an overflow queue.
The overflow queue will concurrently be added to the main queue when a new spot is available.

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

  // Add objects to the queue
  myQueue.Add(1)
  myQueue.Add(2)
  myQueue.Add(3)
  myQueue.Add(4)
  myQueue.Add(5)

  // Get the next object from the queue and remove it
  myObject := myQueue.Next() // 1
  myObject = myQueue.Next() // 2
  myQueue.Next() // 3

  // Wait for the next object to be available
  myQueue.Next(true) // 4

  // Peek at the next object without removing it from the queue
  myQueue.Peek() // 5
  myQueue.Peek(true) // 5

  // Get the queue size
  myQueue.Len() // 1

  go func(){
    // Note: the Wait method will not run until the queue size reaches 0
    for myQueue.Len() != 0 {
      myQueue.Next(true)
    }
  }()

  // Wait for the queue size to reach 0
  myQueue.Wait()
}

```
