package concurrent

import (
	"sync"
)

type SyncQueue struct {
	ml    *sync.RWMutex
	queue *Queue
}

func NewSyncQueue() *SyncQueue {
	return &SyncQueue{new(sync.RWMutex), NewQueue()}
}

func (q *SyncQueue) Push(ele interface{}) {
	q.ml.Lock()
	defer q.ml.Unlock()
	q.queue.Push(ele)
}

func (q *SyncQueue) Pop() interface{} {
	q.ml.Lock()
	defer q.ml.Unlock()
	return q.queue.Pop()
}

func (q *SyncQueue) Front() interface{} {
	q.ml.RLock()
	defer q.ml.RUnlock()
	return q.queue.Front()
}
func (q *SyncQueue) Back() interface{} {
	q.ml.RLock()
	defer q.ml.RUnlock()
	return q.queue.Back()
}

func (q *SyncQueue) Len() int32 {
	q.ml.RLock()
	defer q.ml.RUnlock()
	return q.queue.Len()
}
