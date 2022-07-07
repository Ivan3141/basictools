package concurrent

import (
	"fmt"
	"go.uber.org/atomic"
	"time"
)

type GoPoolDynamic struct {
	pool           *GoPool
	currJobNum     *atomic.Int32 //已经添加的工作数量
	sendJobFinish  chan bool
	finishedJobNum *atomic.Int32 //已经完成的工作数量
	expireSecond   int           //所有工作接收到之后允许最大执行时间(second)
	finished       chan bool
}

func NewGoPoolDynamic(maxWorkNum, expireSecond int) *GoPoolDynamic {
	return &GoPoolDynamic{NewGoPool(maxWorkNum), atomic.NewInt32(0), make(chan bool), atomic.NewInt32(0), expireSecond, make(chan bool)}
}

//AddNewJob仅支持串行调用
func (goPoolDynamic *GoPoolDynamic) AddNewJob(job *Job) bool {
	goPoolDynamic.currJobNum.Add(1)
	f := func() (funcErr error) {
		var jobErr error
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
				goPoolDynamic.finishedJobNum.Add(1)
			}
		}()
		jobErr = job.work()
		return jobErr
	}
	return goPoolDynamic.pool.AddNewJob(NewJob(f))
}

func (goPoolDynamic *GoPoolDynamic) Run() {
	go func() {
		for {
			select {
			case <-goPoolDynamic.sendJobFinish:
				startTime := time.Now()
				go func() {
					for {
						duration := time.Now().Sub(startTime)
						if duration.Seconds() >= float64(goPoolDynamic.expireSecond) || goPoolDynamic.finishedJobNum.Load() == goPoolDynamic.currJobNum.Load() {
							goPoolDynamic.pool.Stop()
							goPoolDynamic.finished <- true
							close(goPoolDynamic.finished)
							return
						}
					}
				}()
				return
			case <-time.After(time.Duration(goPoolDynamic.expireSecond) * time.Second):
				goPoolDynamic.pool.Stop()
				goPoolDynamic.finished <- true
				close(goPoolDynamic.finished)
				return
			}
		}
	}()
	goPoolDynamic.pool.Run()
}

//工作发送完成后需要调用
func (goPoolDynamic *GoPoolDynamic) SendFinish() {
	goPoolDynamic.sendJobFinish <- true
	close(goPoolDynamic.sendJobFinish)
}

//等待执行完,如果所有任务执行成功则返回true，否则返回false
func (goPoolDynamic *GoPoolDynamic) Wait() bool {
	for {
		select {
		case <-goPoolDynamic.finished:
			if goPoolDynamic.finishedJobNum.Load() == goPoolDynamic.currJobNum.Load() {
				return true
			}
			return false
		case <-time.After(time.Duration(goPoolDynamic.expireSecond) * time.Second):
			return false
		}
	}
}
