package concurrent

import (
	"fmt"
	"testing"
)

func Test_Queue(t *testing.T) {
	q := NewQueue()
	q.Push(1)
	q.Push(2)
	q.Push(3)
	fmt.Println(q.Front())
	fmt.Println(q.Pop())
	fmt.Println(q.Front())
	fmt.Println(q.Back())
	fmt.Println(q.Len())
}
