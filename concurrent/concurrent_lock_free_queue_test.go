package concurrent

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_FeeLockQueue(t *testing.T) {
	q := NewLightQueue()
	q.pushHead(1)
	q.pushHead(2)
	q.pushHead(3)

	time.Sleep(time.Duration(2) * time.Second)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		fmt.Println(q.popTail())
	}()
	go func() {
		defer wg.Done()
		fmt.Println(q.popTail())
	}()
	go func() {
		defer wg.Done()
		fmt.Println(q.popTail())
	}()
	wg.Wait()
}
