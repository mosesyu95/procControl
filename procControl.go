package procControl

import (
	"context"
	"sync"
)

type ProcControl struct {
	sem chan struct{}  // 信号量通道
	wg  sync.WaitGroup // 只跟踪已获得许可的任务
}

func NewProcControl(maxConcurrent int) *ProcControl {
	if maxConcurrent < 1 {
		panic("并发数必须 ≥1")
	}
	return &ProcControl{
		sem: make(chan struct{}, maxConcurrent),
	}
}

// Acquire 安全获取许可（原子化操作）
func (p *ProcControl) Acquire(ctx context.Context) error {
	select {
	case p.sem <- struct{}{}: // 1. 先获取信号量
		p.wg.Add(1) // 2. 再增加计数器
		return nil
	case <-ctx.Done():
		return ctx.Err() // 直接返回错误
	}
}

// Release 安全释放许可（支持幂等调用）
func (p *ProcControl) Release() {
	select {
	case <-p.sem: // 有许可可释放
		p.wg.Done()
	default: // 无许可时安全跳过
		// 可添加日志记录异常情况
	}
}

// TryAcquire 非阻塞获取
func (p *ProcControl) TryAcquire() bool {
	select {
	case p.sem <- struct{}{}:
		p.wg.Add(1)
		return true
	default:
		return false
	}
}

// Wait 等待所有已获取许可的任务
func (p *ProcControl) Wait() {
	p.wg.Wait()
}
