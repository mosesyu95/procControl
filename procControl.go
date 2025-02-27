package procControl

import (
	"context"
	"sync"
	"time"
)

type ProcControl struct {
	wg      sync.WaitGroup
	channel chan struct{} // 使用空结构体节省内存
}

// NewProcControl 创建并发控制器，maxConcurrent为最大并发数
func NewProcControl(maxConcurrent int) *ProcControl {
	return &ProcControl{
		channel: make(chan struct{}, maxConcurrent),
	}
}

// Acquire 获取执行权限（带阻塞）
func (p *ProcControl) Acquire() {
	p.channel <- struct{}{} // 阻塞直到有缓冲区
	p.wg.Add(1)
}

// TryAcquire 尝试获取权限（非阻塞，成功返回true）
func (p *ProcControl) TryAcquire() bool {
	select {
	case p.channel <- struct{}{}:
		p.wg.Add(1)
		return true
	default:
		return false
	}
}

// AcquireWithTimeout 获取执行权限（带超时）
func (p *ProcControl) AcquireWithTimeout(timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case p.channel <- struct{}{}: // 尝试获取信号量
		p.wg.Add(1)
		return true
	case <-ctx.Done(): // 超时或取消
		return false
	}
}

// Release 释放执行权限
func (p *ProcControl) Release() {
	<-p.channel
	p.wg.Done()
}

// Wait 等待所有任务完成
func (p *ProcControl) Wait() {
	p.wg.Wait()
	close(p.channel)
}
