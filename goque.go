package goque

import (
	"fmt"
	"time"
)

// 32 bit limit
// const queueSize uintptr = 4294967295

// 16 bit limit
const queueSize uintptr = 65535

// 8 bit limit
// const queueSize uintptr = 255

type Queue[T any] struct {
	data *queueData

	queue *[queueSize+1]T
	overflow *[]T
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

type qVal[T any] struct {
	mode uint8
	val T
}

// New creates a new queue instance
func New[T any]() *Queue[T] {
	start := uint16(0)
	end := uint16(0)

	size := uintptr(0)
	rmSize := uintptr(0)

	queue := [queueSize+1]T{}
	overflow := []T{}

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
				if size > queueSize {
					overflow = append(overflow, inp.val)
					size++
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
	q.in <- qVal[T]{1, value}
}

// Next grabs the next item from the queue, and removes it
func (q *Queue[T]) Next() T {
	if !q.wait() {
		return q.null
	}

	val := q.queue[*q.data.start]
	*q.data.start++
	*q.data.rmSize++

	q.in <- qVal[T]{mode: 2}

	return val
}

// Peek peeks at the next item in the queue without removing it
func (q *Queue[T]) Peek() T {
	if !q.wait() {
		return q.null
	}

	val := q.queue[*q.data.start]
	return val
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
		fmt.Println(q.Len())
		time.Sleep(10 * time.Millisecond)
	}
	q.Stop()
}
