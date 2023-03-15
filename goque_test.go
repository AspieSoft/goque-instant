package goque

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T){
	myq := New[func()]()

	testSize := int(queueSize)*100
	// testSize := int(queueSize)*10
	// testSize := int(queueSize)*3
	// testSize := int(queueSize)+1

	res := [int(queueSize)*100]int{}

	go func(){
		for i := 0; i < testSize; i++ {
			func(i int){
				myq.Add(func() {
					res[i] = i
				})
			}(i)
		}
	}()

	go func(){
		for i := 0; i < testSize; i++ {
			func(i int){
				myq.Add(func() {
					res[i] = i * -1
				})
			}(i)
		}
	}()

	for i := 0; i < testSize*2; i++ {
		fn := myq.Next()
		go fn()
	}

	time.Sleep(100 * time.Millisecond)

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
