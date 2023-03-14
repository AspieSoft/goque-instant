package goque

import (
	"time"
)

const queueSize uint16 = 65535

type Queue[T any] struct {
	start *uint16
	end *uint16

	queue *[uint32(queueSize)+1]qVal[T]
	overflow *[]qVal[T]
	null T

	// shiftOverflow chan bool
	add chan qVal[T]
}

type qVal[T any] struct {
	val T
	hasVal bool
	overflow bool
}

func New[T any]() *Queue[T] {
	start := uint16(0)
	end := uint16(0)

	add := make(chan qVal[T])

	q := &Queue[T]{
		start: &start,
		end: &end,

		queue: &[uint32(queueSize)+1]qVal[T]{},
		overflow: &[]qVal[T]{},

		add: add,
	}

	go func(){
		for {
			val := <-add
			if val.hasVal {
				if val.overflow {
					*q.overflow = append(*q.overflow, val)
				}else{
					(*q.queue)[*q.end] = val
					*q.end++
				}
			}else if val.overflow {
				if len(*q.overflow) != 0 && !(*q.queue)[*q.end].hasVal {
					(*q.queue)[*q.end] = (*q.overflow)[0]
					*q.overflow = (*q.overflow)[1:]
					*q.end++
				}
			}else{
				break
			}
		}
	}()

	return q
}

/* func (q *Queue[T]) wait() func() {
	for *q.running > 0 {
		time.Sleep(10 * time.Nanosecond)
	}
	time.Sleep(1 * time.Nanosecond)
	if *q.running > 1 {
		time.Sleep(10 * time.Nanosecond)
		*q.running--
		time.Sleep(10 * time.Nanosecond)
		return q.wait()
	}

	return func(){
		if *q.running != 0 {
			*q.running--
		}
		time.Sleep(10 * time.Nanosecond)
	}
} */

func (q *Queue[T]) Add(value T){
	if /* q.queue[*q.end].hasVal */ (*q.queue)[*q.end+1].hasVal {
		q.add <- qVal[T]{value, true, true}
		return
	}

	q.add <- qVal[T]{value, true, false}
}

func (q *Queue[T]) Next() T {
	val := q.queue[*q.start]

	if !val.hasVal {
		if len(*q.overflow) == 0 {
			return q.null
		}

		q.add <- qVal[T]{overflow: true}
		time.Sleep(10 * time.Nanosecond)
		return q.Next()
	}

	q.queue[*q.start] = qVal[T]{}
	*q.start++

	q.add <- qVal[T]{overflow: true}
	return val.val
}

func (q *Queue[T]) Empty() bool {
	return !q.queue[*q.start].hasVal && len(*q.overflow) == 0
}

func (q *Queue[T]) Stop() {
	q.add <- qVal[T]{}
}