package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func Test_GoPoolDynamic(t *testing.T) {
	p := NewGoPoolDynamic(10, 3)
	queue := NewSyncQueue()
	for i := 0; i < 20; i++ {
		id := i
		f := func() error {
			if id%2 == 0 {
				time.Sleep(time.Duration(1) * time.Second)
			}
			fmt.Println(id)
			queue.Push(id)
			return nil
		}
		p.AddNewJob(NewJob(f))
	}
	p.Run()
	p.SendFinish()
	succeed := p.Wait()
	fmt.Printf("size of queue is %d, succeed is %v \n,", queue.Len(), succeed)
}
