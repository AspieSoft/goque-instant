package goque

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T){
	myq := New[int]()

	/* for i := 0; i < int(queueSize)+1; i++ {
		myq.Add(i+1)
	} */

	/* myq.Add(1)
	myq.Add(2)
	myq.Add(3)
	myq.Add(4)
	myq.Add(5)

	for i := 0; i < int(queueSize)+1+5; i++ {
		v := myq.Next()
		if v == 0 {
			// fmt.Println("found 0")
		}
	}

	// fmt.Println(*myq.start, *myq.end)

	return */

	//todo: fix large queue breaking (may be skiping some elements when adding them, or possibly getting them)
	// Only an issue in async method

	// testSize := int(queueSize)*100
	testSize := int(queueSize)*10
	// testSize := int(queueSize)*3
	// testSize := int(queueSize)-100
	// testSize := 100


	// done := false
	go func(){
		for i := 0; i < testSize; i++ {
			myq.Add(i+1)
			time.Sleep(10 * time.Nanosecond)
		}
		// done = true
	}()

	//todo: issue only seems to affect concurrent reads with writes, and not concurrent writes directly
	// maybe the overflow could be the issue

	/* for !done {
		time.Sleep(100 * time.Millisecond)
	} */

	// time.Sleep(100 * time.Nanosecond)
	time.Sleep(100 * time.Millisecond)


	res := []int{}
	for i := 0; i < testSize; i++ {
		v := myq.Next()
		if v == 0 {
			time.Sleep(100 * time.Nanosecond)
			continue
		}
		res = append(res, v)
		time.Sleep(10 * time.Nanosecond)
	}

	time.Sleep(100 * time.Millisecond)

	queueSize := 0
	for i := 0; i < len(myq.queue); i++ {
		if myq.queue[i].hasVal || myq.queue[i].val != 0 {
			queueSize++
		}
	}
	fmt.Println("queue size:", queueSize, "\nqueue over:", len(*myq.overflow))

	fmt.Println(*myq.start, *myq.end)

	if len(res) != testSize {
		t.Error("result did not match expected length\nLength", "\nexp:", testSize, "\ngot:", len(res))
		return
	}

	for i := 0; i < len(res); i++ {
		if res[i] != i+1 {
			t.Error("result did not match expected output\nOutput", "\nexp:", i+1, "\ngot:", res[i])
			return
		}
	}
}
