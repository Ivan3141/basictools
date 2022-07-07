package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func Test_GoPool(t *testing.T) {
	pool := NewGoPool(10)
	queue := NewSyncQueue()
	pool.Run()
	for i := 0; i < 100; i++ {
		id := i
		f := func() error {
			queue.Push(id)
			return nil
		}
		pool.AddNewJob(NewJob(f))
	}
	time.Sleep(time.Duration(2) * time.Second)
	pool.Stop()
	fmt.Printf("queue size is %d \n", queue.Len())
}
