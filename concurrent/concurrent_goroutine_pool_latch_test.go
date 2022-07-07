package concurrent

import (
	"fmt"
	"testing"
)

func Test_GoPoolLatch(t *testing.T) {
	p := NewGoPoolLatch(10, 20)
	queue := NewSyncQueue()
	p.Run()
	for i := 0; i < 30; i++ {
		id := i
		f := func() error {
			fmt.Println(id)
			queue.Push(id)
			return nil
		}
		judge := p.AddNewJob(NewJob(f))
		if judge {
			fmt.Printf("id %d is  %d \n", id, 1)
		} else {
			fmt.Printf("id %d is  %d \n", id, 0)
		}
	}
	p.Wait()
	fmt.Println(queue.Len())
}
