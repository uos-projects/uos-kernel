package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/uos-projects/uos-kernel/kernel"
)

// Process 定义一个用户态业务流程
type Process interface {
	// Name 返回流程名称
	Name() string
	// Run 运行流程，阻塞直到 ctx 取消或流程结束
	Run(ctx context.Context, k *kernel.Kernel) error
}

// App 应用框架，管理业务流程的生命周期
type App struct {
	Kernel *kernel.Kernel
}

// New 创建应用
func New(k *kernel.Kernel) *App {
	return &App{Kernel: k}
}

// Run 启动所有业务流程，阻塞直到 ctx 取消或所有流程结束
func (a *App) Run(ctx context.Context, processes ...Process) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(processes))

	for _, p := range processes {
		wg.Add(1)
		go func(proc Process) {
			defer wg.Done()
			if err := proc.Run(ctx, a.Kernel); err != nil {
				errCh <- fmt.Errorf("process %s: %w", proc.Name(), err)
			}
		}(p)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("app errors: %v", errs)
	}
	return nil
}
