package goque

import (
	"time"
)

// const queueSize uint16 = 65535
const queueSize uint8 = 255

type Queue[T any] struct {
	start *uint8
	end *uint8
	size *uintptr
	rmSize *uintptr

	queue *[uint32(queueSize)+1]qVal[T]
	overflow *[]qVal[T]
	null T

	in chan qVal[T]
	fixing *bool
}

type qVal[T any] struct {
	mode uint8
	val T
}

func New[T any]() *Queue[T] {
	start := uint8(0)
	end := uint8(0)
	size := uintptr(0)
	rmSize := uintptr(0)

	queue := [uint32(queueSize)+1]qVal[T]{}
	overflow := []qVal[T]{}

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
				if size > uintptr(queueSize) {
					overflow = append(overflow, inp)
					size++
				}else{
					queue[end] = inp
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
			}
		}
	}()

	return &q
}

// wait will try to wait and check again if the queue reports empty
//
// @return bool
//
// true = finished
//
// false = timeout
func (q *Queue[T]) wait(value T) bool {
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

func (q *Queue[T]) Add(value T){
	q.in <- qVal[T]{1, value}
}

func (q *Queue[T]) Next() T {
	// if empty, try to wait and check again
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
			return q.null
		}
	}

	val := q.queue[*q.start]
	*q.start++
	*q.rmSize++

	q.in <- qVal[T]{mode: 2}

	return val.val
}

func (q *Queue[T]) Len() uintptr {
	// if empty, try to wait and check again
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
	}

	return *q.size - *q.rmSize
}
