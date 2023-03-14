# Goque Instant

[![donation link](https://img.shields.io/badge/buy%20me%20a%20coffee-square-blue)](https://buymeacoffee.aspiesoft.com)

A fast queue system for golang.

This module avoids shifting values to deferent indexes of an array to different spots in memory.

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

  if myQueue.Empty() {
    myQueue.Stop()
  }
}

```
