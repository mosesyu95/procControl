package procControl

import (
	"sync"
)

type ProcControl struct {
	wg      *sync.WaitGroup
	channel chan int8
}

func NewProcControl(procNum int) *ProcControl {
	return &ProcControl{
		wg:      &sync.WaitGroup{},
		channel: make(chan int8, procNum),
	}
}

func (p *ProcControl) Acquire() {
	p.wg.Add(1)
	p.channel <- int8(0)
}

func (p *ProcControl) Release() {
	p.wg.Done()
	<-p.channel
}

func (p *ProcControl) Wait() {
	p.wg.Wait()
}
