package concurrent

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

//noop is implemented by javen.you, here is the version with little change on initialization of new noop
type Noop struct {
	startWork    chan bool
	stopWork     chan bool
	wg           sync.WaitGroup
	stopFlag     bool
	signalCh     chan os.Signal
	sleepTimeSec int
	fc           func()
}

func NewNoop(sleepTime int, fnc func()) *Noop {
	return &Noop{make(chan bool), make(chan bool), sync.WaitGroup{}, false, make(chan os.Signal), sleepTime, fnc}

}

func (noop *Noop) Run() {

	// 监听关闭信号，做平滑重启使用
	signal.Notify(noop.signalCh, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)
	go func() {
		sig := <-noop.signalCh
		switch sig {
		default:
		case syscall.SIGHUP:
		case syscall.SIGINT:
			noop.stopWork <- true
		case syscall.SIGUSR1:
		case syscall.SIGUSR2:
		}
	}()
	noop.wg.Add(1)
	go func() {
		defer noop.wg.Done()
	Loop:
		for {
			select {
			case <-noop.startWork:
				go func() {
					noop.wg.Add(1)
					defer noop.wg.Done()
					noop.fc()
					if !noop.stopFlag {
						time.Sleep(time.Second * time.Duration(noop.sleepTimeSec))
						noop.startWork <- true
					}
				}()
			case <-noop.stopWork:
				noop.stopFlag = true
				break Loop
			}
		}
	}()
	// 第一次启动程task
	noop.startWork <- true
	noop.wg.Wait()
}
