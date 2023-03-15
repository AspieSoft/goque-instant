# Goque Instant

[![donation link](https://img.shields.io/badge/buy%20me%20a%20coffee-paypal-blue)](https://paypal.me/shaynejrtaylor?country.x=US&locale.x=en_US)

A fast queue system for golang.

This module avoids shifting values to different indexes of an array and to different spots in memory.

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

  if myQueue.Len() == 0 {
    myQueue.Stop()
  }
}

```
