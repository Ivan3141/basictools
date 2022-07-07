package concurrent

import (
	"container/list"
)

type Queue struct {
	queue *list.List
}

func NewQueue() *Queue {
	return &Queue{list.New()}
}

func (q *Queue) Push(ele interface{}) {
	q.queue.PushBack(ele)
}

func (q *Queue) Pop() interface{} {
	ele := q.queue.Front()
	if ele != nil {
		val := q.queue.Remove(ele)
		return val
	}
	return nil
}

func (q *Queue) Front() interface{} {
	ele := q.queue.Front()
	if ele != nil {
		return ele.Value
	}
	return nil
}
func (q *Queue) Back() interface{} {
	ele := q.queue.Back()
	if ele != nil {
		return ele.Value
	}
	return nil
}

func (q *Queue) Len() int32 {
	return int32(q.queue.Len())
}
