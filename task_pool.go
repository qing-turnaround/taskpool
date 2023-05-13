package pool

import (
	"context"
	"fmt"
)

// 选项模式：用于设置任务池的属性
type PoolOptions struct {
	// 任务池的大小（启动多少goroutine）
	poolSize int
	// 任务池的缓冲区大小（设置taskPool能存放多少任务）
	bufferSize int
}

// 选项模式：用于设置任务池的缓冲区大小
func WithBufferSize(bufferSize int) func(*PoolOptions) {
	return func(options *PoolOptions) {
		options.bufferSize = bufferSize
	}
}

// 选项模式：用于设置任务池的 Goroutine 数量
func WithPoolSize(poolSize int) func(*PoolOptions) {
	return func(options *PoolOptions) {
		options.poolSize = poolSize
	}
}

// task 定义任务
type task func()

func (t task) Run() {
	// 捕获panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("task panic:", err)
		}

	}()
	// 执行任务
	t()
}

type TaskPool struct {
	// 任务队列
	taskChan chan task
	// 任务池的关闭通道
	closeChan chan struct{}
}

// NewTaskPool 创建一个任务池
func NewTaskPool(options ...func(*PoolOptions)) *TaskPool {
	// 设置默认值
	poolOptions := &PoolOptions{
		poolSize:   10,
		bufferSize: 100,
	}
	// 通过选项模式设置任务池的属性
	for _, option := range options {
		option(poolOptions)
	}
	// 创建任务池
	taskPool := &TaskPool{
		taskChan:  make(chan task, poolOptions.bufferSize),
		closeChan: make(chan struct{}),
	}

	// 启动任务池
	for i := 0; i < poolOptions.poolSize; i++ {
		go func() {
			for {
				select {
				case task := <-taskPool.taskChan:
					task.Run()
				case <-taskPool.closeChan:
					return
				}
			}
		}()
	}

	return taskPool
}

// AddTask 添加任务-可以通过context控制 添加任务的超时
func (t *TaskPool) AddTask(task task, ctx context.Context) error {
	// 将任务添加到任务队列中
	select {
	case t.taskChan <- task:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// Close 关闭任务池
// 请不要重复调用，否则会panic（也可以使用sync.Once，但有时候 我们不应该过度保姆式 设计）
func (t *TaskPool) Close() {
	close(t.closeChan)
}
