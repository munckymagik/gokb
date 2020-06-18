package scheduler

import "time"

type donetoken struct{}

type Scheduler struct {
	interval time.Duration
	onTick   func(bool)
	stopChan chan chan<- donetoken
}

func New(interval time.Duration, onTick func(bool)) Scheduler {
	return Scheduler{interval, onTick, nil}
}

func (scheduler *Scheduler) Start() {
	stopChan := make(chan chan<- donetoken, 1)
	go func() {
		ticker := time.NewTicker(scheduler.interval)
		defer ticker.Stop()
		scheduler.loopWithTicker(ticker.C, stopChan)
	}()

	scheduler.stopChan = stopChan
}

func (scheduler *Scheduler) Stop() {
	done := make(chan donetoken, 1)
	scheduler.stopChan <- done
	<-done
}

func (scheduler *Scheduler) loopWithTicker(tickSource <-chan time.Time, stop chan chan<- donetoken) {
	for {
		select {
		case <-tickSource:
			scheduler.onTick(false)
		case done := <-stop:
			scheduler.onTick(true)
			done <- donetoken{}
			return
		}
	}
}
