package goque

import (
	"fmt"
	"math"
	"time"

	"github.com/pbnjay/memory"
)

// 32 bit limit
// const queueSize uintptr = 4294967295+1

// 16 bit limit
const queueSize uintptr = 65535+1

// 8 bit limit
// const queueSize uintptr = 255+1

type Queue[T any] struct {
	data *queueData

	queue *[queueSize]qObj[T]
	overflow *[]qObj[T]
	null T

	in chan qVal[T]
}

type queueData struct {
	start *uint16
	end *uint16
	size *uintptr
	rmSize *uintptr
	fixing *bool
}

type qObj[T any] struct {
	val T
	hasVal bool
}

type qVal[T any] struct {
	mode uint8
	val qObj[T]
	start uint16
}

var memoryUsageAvailable float64

func init(){
	go func(){
		for {
			memoryUsageAvailable = formatMemoryUsage(memory.FreeMemory())

			if memoryUsageAvailable > 10000 { // 10gb
				time.Sleep(10 * time.Millisecond)
			} else if memoryUsageAvailable > 2000 { // 2gb
				time.Sleep(10000 * time.Nanosecond)
			} else if memoryUsageAvailable > 1000 { // 1gb
				time.Sleep(100 * time.Nanosecond)
			}else if memoryUsageAvailable > 500 { // 500mb
				time.Sleep(10 * time.Nanosecond)
			}else if memoryUsageAvailable > 250 { // 250mb
				time.Sleep(1 * time.Nanosecond)
			}else{
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

var reportedLowMem bool = false

// ConsoleLogsEnabled can be modified to change logging rules for this module
//
// 0 = disabled
//
// 1 = normal
//
// 2 = debug
var ConsoleLogsEnabled uint8 = 1

func reportLowMem(){
	if !reportedLowMem {
		reportedLowMem = true
		if ConsoleLogsEnabled >= 1 {
			fmt.Println("Low Memory Detected:", memoryUsageAvailable, "[goque-instant]")
		}
	}
}

func reportStableMem(){
	if reportedLowMem {
		reportedLowMem = false
		if ConsoleLogsEnabled >= 1 {
			fmt.Println("Stable Memory Recovered:", memoryUsageAvailable, "[goque-instant]")
		}
	}
}

// New creates a new queue instance
func New[T any]() *Queue[T] {
	start := uint16(0)
	end := uint16(0)

	size := uintptr(0)
	rmSize := uintptr(0)

	queue := [queueSize]qObj[T]{}
	overflow := []qObj[T]{}
	null := qObj[T]{}
	
	in := make(chan qVal[T])
	fixing := false

	q := Queue[T]{
		data: &queueData{
			start: &start,
			end: &end,
			size: &size,
			rmSize: &rmSize,
			fixing: &fixing,
		},

		queue: &queue,
		overflow: &overflow,

		in: in,
	}

	go func(){
		for {
			inp := <-in
			if inp.mode == 1 { // Add
				if size >= queueSize {
					overflow = append(overflow, inp.val)
					size++

					// reduce high memory usage
					// if size > 24000000 || (size > 4800000 && memoryUsageAvailable < 500) {
					if size > 4800000 || (size > 2400000 && memoryUsageAvailable < 500) {
						time.Sleep(100 * time.Nanosecond)

						if memoryUsageAvailable < 500 && memoryUsageAvailable != 0 {
							if memoryUsageAvailable < 250 {
								reportLowMem()
							}

							loops := int(10000 - math.Floor((memoryUsageAvailable * 4500 / 300)))
							time.Sleep(time.Duration(loops/10) * time.Nanosecond)

							if memoryUsageAvailable < 300 || size > 24000000 {
								timeout := time.Duration(math.Max(math.Min(float64(size / 10000), 5000), 1000))

								if ConsoleLogsEnabled >= 2 && size % 1000 == 0 {
									fmt.Println("low mem:", memoryUsageAvailable, "size:", size, "timeout:", loops, "-", timeout.Nanoseconds())
								}

								for (memoryUsageAvailable < 300 || size > 24000000) && loops > 0 {
									loops--
									time.Sleep(timeout * time.Nanosecond)
								}
							}

						}
					}else if memoryUsageAvailable > 1000 {
						reportStableMem()
					}
				}else{
					queue[end] = inp.val
					end++
					size++
				}
			}else if inp.mode == 2 { // Next
				if len(overflow) != 0 {
					queue[end] = overflow[0]
					overflow = overflow[1:]
					end++
				}

				fixing = true
				time.Sleep(10 * time.Nanosecond)
				size--
				rmSize--
				if !queue[inp.start].hasVal {
					queue[inp.start] = null
				}
				fixing = false
			}else{ // Stop
				break
			}
		}
	}()

	return &q
}

// wait will try to wait and check again if the queue reports empty, but may still be adding items
//
// @return bool
//
// true = finished
//
// false = timeout
func (q *Queue[T]) wait() bool {
	for *q.data.fixing {
		time.Sleep(10 * time.Nanosecond)
	}

	if *q.data.size - *q.data.rmSize == 0 {
		loops := 100000
		for *q.data.fixing || (*q.data.size - *q.data.rmSize == 0 && *q.data.start == *q.data.end && loops > 0) {
			if !*q.data.fixing {
				loops--
			}
			time.Sleep(10 * time.Nanosecond)
		}

		if *q.data.size - *q.data.rmSize == 0 {
			return false
		}
	}

	return true
}

// Add adds an item to the queue
func (q *Queue[T]) Add(value T){
	q.in <- qVal[T]{mode: 1, val: qObj[T]{value, true}}
}

// Next grabs the next item from the queue, and removes it
func (q *Queue[T]) Next() T {
	if !q.wait() {
		return q.null
	}

	val := q.queue[*q.data.start]
	q.queue[*q.data.start].hasVal = false
	*q.data.start++
	*q.data.rmSize++

	q.in <- qVal[T]{mode: 2, start: (*q.data.start)-1}

	return val.val
}

// Peek peeks at the next item in the queue without removing it
func (q *Queue[T]) Peek() T {
	if !q.wait() {
		return q.null
	}

	val := q.queue[*q.data.start]
	return val.val
}

// Len returns the number of items in the queue
func (q *Queue[T]) Len() uintptr {
	q.wait()

	return *q.data.size - *q.data.rmSize
}

// Stop stops the queue from running the loop that adds items concurrently through a channel
func (q *Queue[T]) Stop() {
	q.wait()

	q.in <- qVal[T]{}
}

// Wait waits for the queue to have 0 items left
func (q *Queue[T]) Wait() {
	for q.Len() != 0 {
		time.Sleep(10 * time.Millisecond)
	}
}

// WaitAndStop waits for the queue to have 0 items left, then runs Stop
func (q *Queue[T]) WaitAndStop() {
	for q.Len() != 0 {
		time.Sleep(10 * time.Millisecond)
	}
	q.in <- qVal[T]{}
}

// LowMem returns true if the module detected the system drop below 250mb,
// that may have been caused by a large queue size (> 2400000 objects)
//
// once low memory detection is triggered, it will only return as back to stable if memory usage goes back up to 1gb,
// and if the queue size has dipped down to < 4800000 objects
func LowMem() bool {
	return reportedLowMem
}


func formatMemoryUsage(b uint64) float64 {
	return math.Round(float64(b) / 1024 / 1024 * 100) / 100
}
