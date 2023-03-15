package goque

import (
	"testing"
)

//? 3.5, max: 4.1

func Test(t *testing.T){
	myq := New[int]()

	testSize := int(queueSize)*100
	// testSize := int(queueSize)*10
	// testSize := int(queueSize)*3
	// testSize := int(queueSize)+1

	go func(){
		for i := 0; i < testSize; i++ {
			myq.Add(i)
		}
	}()

	for i := 0; i < testSize; i++ {
		v := myq.Next()

		if v != i {
			t.Error("invalid value:", v, "!=", i)
			break
		}
	}
}
