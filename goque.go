package goque

import (
	"time"
)

const queueSize uint16 = 65535

type Queue[T any] struct {
	start *uint16
	end *uint16
	size *uintptr

	queue *[uint32(queueSize)+1]qVal[T]
	overflow *[]qVal[T]
	null T

	// shiftOverflow chan bool
	add chan qVal[T]
}

type qVal[T any] struct {
	mode uint8
	val T
	// hasVal bool
	// overflow bool
}

func New[T any]() *Queue[T] {
	start := uint16(0)
	end := uint16(0)
	size := uintptr(0)

	queue := [uint32(queueSize)+1]qVal[T]{}
	overflow := []qVal[T]{}

	add := make(chan qVal[T])

	q := &Queue[T]{
		start: &start,
		end: &end,
		size: &size,

		queue: &queue,
		overflow: &overflow,

		add: add,
	}

	go func(){
		for {
			val := <-add
			if val.mode == 1 {
				if size > uintptr(queueSize) {
					overflow = append(overflow, val)
				}else{
					queue[end] = val
					size++
					end++
				}
			}else if val.mode >= 2 {
				if val.mode == 2 && size > 0 {
					size--
				}

				if len(overflow) != 0 && size <= uintptr(queueSize) {
					queue[end] = overflow[0]
					overflow = overflow[1:]
					end++
				}
			}else{
				break
			}
		}
	}()

	return q
}

func (q *Queue[T]) Add(value T){
	q.add <- qVal[T]{mode: 1, val: value}
}

func (q *Queue[T]) Next() T {
	if *q.size == 0 {
		return q.null
	}

	val := q.queue[*q.start]

	for val.mode == 0 {
		q.add <- qVal[T]{mode: 3}
		time.Sleep(10 * time.Nanosecond)
		val = q.queue[*q.start]

		if val.mode == 0 && *q.size == 0 {
			return q.null
		}
	}

	*q.start++
	q.add <- qVal[T]{mode: 2}
	return val.val
}

func (q *Queue[T]) Empty() bool {
	return *q.size == 0
}

func (q *Queue[T]) Stop() {
	q.add <- qVal[T]{}
}
