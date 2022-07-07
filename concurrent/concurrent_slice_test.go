package concurrent

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_SyncSlice(t *testing.T) {
	s := NewSyncSlice()
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		s.Append([]int{1, 3, 4})
	}()
	go func() {
		defer wg.Done()
		s.Append([]int{4})
	}()
	go func() {
		defer wg.Done()
		s.Append([]int{9, 10})
	}()
	wg.Wait()
	time.Sleep(time.Duration(2) * time.Second)
	slices := s.GetSlice().([]interface{})
	for _, slice := range slices {
		fmt.Println(slice.([]int))
	}
}
