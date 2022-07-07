package concurrent

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_SyncQueue(t *testing.T) {
	q := NewSyncQueue()

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		q.Push(1)
	}()
	go func() {
		defer wg.Done()
		q.Push(2)
	}()
	go func() {
		defer wg.Done()
		q.Push(3)
	}()
	wg.Wait()
	time.Sleep(time.Duration(2) * time.Second)
	fmt.Println(q.Len())
	fmt.Println(q.Front())
	fmt.Println(q.Back())
	wg.Add(3)
	go func() {
		defer wg.Done()
		q.Push(4)
	}()
	go func() {
		defer wg.Done()
		fmt.Println(q.Pop())
	}()
	go func() {
		defer wg.Done()
		fmt.Println(q.Pop())
	}()
	wg.Wait()
	time.Sleep(time.Duration(2) * time.Second)
	fmt.Println(q.Front())
	fmt.Println(q.Back())
	fmt.Println(q.Len())
}
