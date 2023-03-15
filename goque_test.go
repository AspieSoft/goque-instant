package goque

import (
	"testing"
	"time"
)

func Test(t *testing.T){
	myq := New[int]()

	// testSize := int(queueSize)*100
	// testSize := int(queueSize)*10
	testSize := int(queueSize)*3
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
		res = append(res, v)
		time.Sleep(10 * time.Nanosecond)
	}

	time.Sleep(100 * time.Millisecond)

	// fmt.Println(myq.overflow)

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
