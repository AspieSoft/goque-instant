package goque

import (
	"time"
)

// 32 bit limit
// const queueSize uintptr = 4294967295

// 16 bit limit
const queueSize uintptr = 65535

//? test1: avr: 4.5s, min: 4.2s, max: 4.9s

// 8 bit limit
// const queueSize uint8 = 255

type queuePos struct {
	mode uint8

	start16 *uint16
	end16 *uint16

	start8 *uint16
	end8 *uint16
}

type Queue[T any] struct {
	start *uint16
	end *uint16
	size *uintptr
	rmSize *uintptr

	queue *[queueSize+1]T
	overflow *[]T
	null T

	in chan qVal[T]
	fixing *bool
}

type qVal[T any] struct {
	mode uint8
	val T
}

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
		start: &start,
		end: &end,
		size: &size,
		rmSize: &rmSize,

		queue: &queue,
		overflow: &overflow,

		in: in,
		fixing: &fixing,
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
	for *q.fixing {
		time.Sleep(10 * time.Nanosecond)
	}

	if *q.size - *q.rmSize == 0 {
		loops := 100000
		for *q.fixing || (*q.size - *q.rmSize == 0 && *q.start == *q.end && loops > 0) {
			if !*q.fixing {
				loops--
			}
			time.Sleep(10 * time.Nanosecond)
		}

		if *q.size - *q.rmSize == 0 {
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

	val := q.queue[*q.start]
	*q.start++
	*q.rmSize++

	q.in <- qVal[T]{mode: 2}

	return val
}

// Peek peeks at the next item in the queue without removing it
func (q *Queue[T]) Peek() T {
	if !q.wait() {
		return q.null
	}

	val := q.queue[*q.start]
	return val
}

func (q *Queue[T]) Len() uintptr {
	q.wait()

	return *q.size - *q.rmSize
}

func (q *Queue[T]) Stop() {
	q.wait()

	q.in <- qVal[T]{}
}


func (q *queuePos) getStart() uint {
	switch q.mode {
	case 16:
		return uint(*q.start16)
	case 8:
		return uint(*q.start8)
	default:
		return 0
	}
}

func (q *queuePos) getEnd() uint {
	switch q.mode {
	case 16:
		return uint(*q.end16)
	case 8:
		return uint(*q.end8)
	default:
		return 0
	}
}

func (q *queuePos) addStart() {
	switch q.mode {
	case 16:
		*q.start16++
	case 8:
		*q.start8++
	}
}

func (q *queuePos) addEnd() {
	switch q.mode {
	case 16:
		*q.end16++
	case 8:
		*q.end8++
	}
}
