package procControl

import (
	"sync"
)

type ProcControl struct {
	wg      sync.WaitGroup
	channel chan int8
}

func NewProcControl(procNum, totalProc int) *ProcControl {
	wg := sync.WaitGroup{}
	wg.Add(totalProc)
	return &ProcControl{
		wg:      wg,
		channel: make(chan int8, procNum),
	}
}

func (p *ProcControl) Acquire() {
	p.channel <- int8(0)
}

func (p *ProcControl) Release() {
	p.wg.Done()
	<-p.channel
}

func (p *ProcControl) Done() {
	p.wg.Done()
}

func (p *ProcControl) Wait() {
	p.wg.Wait()
}
