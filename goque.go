package goque

import (
	"sync"
	"time"
)

// 32 bit limit
// const queueSize uintptr = 4294967295+1

// 16 bit limit
const queueSize uintptr = 65535+1

// 8 bit limit
// const queueSize uintptr = 255+1

type Queue[T any] struct {
	mu sync.RWMutex

	queue *[queueSize]T
	overflow *[]T
	null T

	data *qData
}

type qData struct {
	mu sync.Mutex

	start *uint16
	end *uint16
	size *uintptr
}

// edit incraments a number by 1
//
// 's' = start++
//
// 'e' = end++
//
// '+' = size++
//
// '-' = size--
func (q *qData) edit(n ...byte){
	q.mu.Lock()
	for _, v := range n {
		switch v {
		case 's':
			*q.start++
		case 'e':
			*q.end++
		case '+':
			*q.size++
		case '-':
			*q.size--
		}
	}
	q.mu.Unlock()
}

func New[T any]() *Queue[T] {
	start := uint16(0)
	end := uint16(0)
	size := uintptr(0)

	return &Queue[T]{
		queue: &[queueSize]T{},
		overflow: &[]T{},

		data: &qData{
			start: &start,
			end: &end,
			size: &size,
		},
	}
}

func (q *Queue[T]) Add(value T){
	q.mu.Lock()
	defer q.mu.Unlock()

	if *q.data.size >= queueSize {
		*q.overflow = append(*q.overflow, value)
		*q.data.size++
		return
	}

	q.queue[*q.data.end] = value
	q.data.edit('e', '+')
}

func (q *Queue[T]) Next(wait ...bool) T {
	q.mu.RLock()
	defer q.mu.RUnlock()

	loops := 10000
	for *q.data.size == 0 && loops > 0 {
		if len(wait) == 0 || wait[0] == false {
			loops--
		}
		q.mu.RUnlock()
		time.Sleep(100 * time.Nanosecond)
		q.mu.RLock()
	}

	if *q.data.size == 0 {
		return q.null
	}

	start := *q.data.start
	q.data.edit('s', '-')

	val := q.queue[start]
	q.queue[start] = q.null

	go func(){
		q.mu.Lock()
		defer q.mu.Unlock()

		if len(*q.overflow) != 0 {
			q.queue[*q.data.end] = (*q.overflow)[0]
			*q.overflow = (*q.overflow)[1:]
			q.data.edit('e')
		}
	}()

	return val
}

func (q *Queue[T]) Peek(wait ...bool) T {
	q.mu.RLock()
	defer q.mu.RUnlock()

	loops := 10000
	for *q.data.size == 0 && loops > 0 {
		if len(wait) == 0 || wait[0] == false {
			loops--
		}
		q.mu.RUnlock()
		time.Sleep(100 * time.Nanosecond)
		q.mu.RLock()
	}

	if *q.data.size == 0 {
		return q.null
	}

	return q.queue[*q.data.start]
}

func (q *Queue[T]) Len() uintptr {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return *q.data.size
}

func (q *Queue[T]) Wait() {
	time.Sleep(10 * time.Millisecond)

	for {
		time.Sleep(100 * time.Nanosecond)

		if q.mu.TryRLock() {
			if *q.data.size == 0 {
				q.mu.RUnlock()
				break
			}
			q.mu.RUnlock()
		}
	}
}
