package concurrent

import (
	"fmt"
	"go.uber.org/atomic"
	"sync"
)

type GoPoolLatch struct {
	pool        *GoPool
	wg          sync.WaitGroup
	currJobNum  *atomic.Int32 //当前已经执行过的job数量
	totalJobNum int32
}

//maxWorkNum 协程池开启协程数量，totalJobNum要执行的任务的数量
func NewGoPoolLatch(maxWorkNum, totalJobNum int) *GoPoolLatch {
	p := &GoPoolLatch{pool: NewGoPool(int(maxWorkNum)), wg: sync.WaitGroup{}, currJobNum: atomic.NewInt32(0), totalJobNum: int32(totalJobNum)}
	p.wg.Add(totalJobNum)
	return p
}

//AddNewJob 只支持串行调用
func (goPoolLatch *GoPoolLatch) AddNewJob(job *Job) bool {
	jobNum := goPoolLatch.currJobNum.Load()
	if jobNum == goPoolLatch.totalJobNum {
		return false
	}
	goPoolLatch.currJobNum.Add(1)
	newJob := NewJob(func() (funcErr error) {
		var jobErr error
		defer goPoolLatch.wg.Done()
		defer func() {
			if rec := recover(); rec != nil {
				err, isError := rec.(error)
				if !isError {
					err = fmt.Errorf("%v", rec)
				}
				funcErr = err
			} else if jobErr != nil {
				funcErr = jobErr
			} else {
				funcErr = nil
			}
		}()
		jobErr = job.work()
		return jobErr
	})
	return goPoolLatch.pool.AddNewJob(newJob)
}

func (goPoolLatch *GoPoolLatch) Run() {
	goPoolLatch.pool.Run()
}

func (goPoolLatch *GoPoolLatch) Wait() {
	goPoolLatch.wg.Wait()
	goPoolLatch.pool.Stop()
}
