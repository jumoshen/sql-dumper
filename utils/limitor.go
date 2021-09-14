package utils

import (
	"sync"
)

type WaitGroup struct {
	workChan chan int
	wg       sync.WaitGroup
}

func NewPool(limit int) *WaitGroup {
	ch := make(chan int, limit)

	return &WaitGroup{
		workChan: ch,
		wg:       sync.WaitGroup{},
	}
}

func (wg *WaitGroup) Add(num int) {
	for i := 0; i < num; i++ {
		wg.workChan <- i
		wg.wg.Add(1)
	}
}

func (wg *WaitGroup) Done() {
LOOP:
	for {
		select {
		case <-wg.workChan:
			break LOOP
		}
	}
	wg.wg.Done()
}

func (wg *WaitGroup) wait() {
	wg.wg.Wait()
}
