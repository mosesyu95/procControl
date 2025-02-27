package procControl

import (
	"context"
	"sync"
)

type ProcControl struct {
	wg      sync.WaitGroup
	channel chan struct{} // 使用空结构体节省内存
}

// NewProcControl 创建并发控制器，maxConcurrent为最大并发数
func NewProcControl(maxConcurrent int) *ProcControl {
	if maxConcurrent < 1 {
		panic("并发数必须 ≥1")
	}
	return &ProcControl{
		channel: make(chan struct{}, maxConcurrent),
	}
}

// Acquire 获取执行权限（带阻塞）
func (p *ProcControl) Acquire() {
	// 1. 先保证计数器增加
	p.wg.Add(1)

	// 2. 添加 panic 恢复机制
	defer func() {
		if r := recover(); r != nil {
			// 发生 panic 时回滚计数器
			p.wg.Done()
			// 重新抛出 panic（可选）
			panic(r)
		}
	}()

	// 3. 执行信号量操作（可能引发 panic 的位置）
	p.channel <- struct{}{}
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

// AcquireWithTimeout 获取执行权限（可带超时 context.WithTimeout）
func (p *ProcControl) AcquireWithTimeout(ctx context.Context) bool {

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
	select {
	case <-p.channel: // 非阻塞接收
		p.wg.Done()
	default:
		//log.Println("警告：未匹配的Release调用")
		panic("Release called without Acquire")
	}
}

// Wait 等待所有任务完成
func (p *ProcControl) Wait() {
	p.wg.Wait()
	close(p.channel)
}
