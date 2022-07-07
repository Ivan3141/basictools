package concurrent

import (
	"fmt"
	"go.uber.org/atomic"
)

type Job struct {
	f func() error
}

func NewJob(f func() error) *Job {
	t := Job{
		f: f,
	}
	return &t
}

func (t *Job) work() error {
	return t.f()
}

type GoPool struct {
	workerChannel chan *Job
	stopChannel   chan bool
	maxWorkerNum  int           //最大协程数量，由人工设置
	currWorkerNum *atomic.Int32 //当前实时协程数量
	syncQueue     *SyncQueue
	isStop        *atomic.Bool
}

//创建一个协程池
func NewGoPool(maxNum int) *GoPool {
	p := GoPool{
		maxWorkerNum:  maxNum,
		workerChannel: make(chan *Job),
		stopChannel:   make(chan bool),
		currWorkerNum: atomic.NewInt32(0),
		syncQueue:     NewSyncQueue(),
		isStop:        atomic.NewBool(false),
	}
	return &p
}

//AddNewJob仅支持串行调用
func (p *GoPool) AddNewJob(job *Job) bool {
	if p.isStop.Load() {
		return false
	}
	p.syncQueue.Push(job)
	return true
}

func (p *GoPool) Stop() {
	p.stopChannel <- true
	close(p.stopChannel)
}

func (p *GoPool) Run() {
	go func() {
		for {
			select {
			case <-p.stopChannel:
				p.isStop.Store(true)
				close(p.workerChannel)
				return
			default:
				job := p.syncQueue.Pop()
				if job != nil {
					p.workerChannel <- job.(*Job)
				}
			}
		}
	}()
	for i := 0; i < p.maxWorkerNum; i++ {
		go func() {
			p.work()
		}()
	}
}

func (p *GoPool) work() (jobErr error) {
	defer func() {
		if rec := recover(); rec != nil {
			err, isError := rec.(error)
			if !isError {
				err = fmt.Errorf("%v", rec)
			}
			jobErr = err
		}
	}()
	p.currWorkerNum.Add(1)
	for job := range p.workerChannel {
		err := job.work()
		if err != nil {
			return err
		}
	}
	return nil
}
