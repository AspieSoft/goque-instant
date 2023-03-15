package goque

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T){
	ConsoleLogsEnabled = 2

	// const testSize int = int(queueSize)*500
	// const testSize int = int(queueSize)*256
	// const testSize int = int(queueSize)*100
	const testSize int = int(queueSize)*10
	// const testSize int = int(queueSize)*3
	// const testSize int = int(queueSize)+1

	fmt.Println("Test Size:", testSize)

	myq := New[func()]()

	res := [testSize*2]int{}

	go func(){
		for i := 0; i < testSize; i++ {
			func(i int){
				myq.Add(func() {
					res[i] = i
				})
			}(i)
			time.Sleep(10 * time.Nanosecond)
		}
	}()

	go func(){
		for i := 0; i < testSize; i++ {
			func(i int){
				myq.Add(func() {
					res[i] = i * -1
				})
			}(i)
			time.Sleep(10 * time.Nanosecond)
		}
	}()

	go func(){
		for i := 0; i < testSize*2; i++ {
			fn := myq.Next()
			if fn != nil {
				go fn()
			}
			time.Sleep(10 * time.Nanosecond)
		}
	}()

	myq.Wait()

	hasPos := false
	hasNeg := false

	for i := 0; i < testSize; i++ {
		if res[i] != i && res[i] != i * -1 {
			t.Error("invalid value:", res[i], "!=", i)
			break
		}

		if res[i] == i {
			hasPos = true
		}

		if res[i] == i * -1 {
			hasNeg = true
		}
	}

	fmt.Println(hasPos, hasNeg)

	myq.WaitAndStop()
}
