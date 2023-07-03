package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
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
type task func() error

func (t task) Run() (errRet error) {
	// 捕获panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("task panic:", err)
			errRet = errors.New("task panic")
		}
	}()
	// 执行任务
	if err := t(); err != nil {
		return err
	}
	return nil
}

type TaskPool struct {
	// 当前任务池得任务数量
	count int
	// 任务执行失败，或者超时未执行得任务数量
	fails   int
	mu      sync.RWMutex
	options *PoolOptions
	// 任务队列
	taskChan chan task
	// 任务池的关闭通道
	closeChan chan struct{}
}

// NewTaskPool 创建一个任务池
func NewTaskPool(options ...func(*PoolOptions)) *TaskPool {

	poolOptions := &PoolOptions{}
	// 通过选项模式设置任务池的属性
	for _, option := range options {
		option(poolOptions)
	}

	// 创建任务池
	taskPool := &TaskPool{
		taskChan:  make(chan task, poolOptions.bufferSize),
		closeChan: make(chan struct{}),
		mu:        sync.RWMutex{},
		options:   poolOptions,
	}

	// 启动任务池
	for i := 0; i < taskPool.options.poolSize; i++ {
		go func() {
			for {
				select {
				case task := <-taskPool.taskChan:
					if err := task.Run(); err != nil {
						// 执行失败，Fails++
						taskPool.AddFails()
					}
					// 执行完毕，Count--
					taskPool.SubCount()
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
	t.AddCount()
	// 将任务添加到任务队列中
	select {
	// select会阻塞， 直到有一个case完成操作
	case t.taskChan <- task:
	case <-ctx.Done():
		// 超时未执行，Fails++
		t.AddFails()
		t.SubCount()
		return ctx.Err()
		// default:
		// 	return fmt.Errorf("任务池已满")
	}
	return nil
}

// AddCount 增加任务池的任务数量
func (t *TaskPool) AddCount() {
	defer t.mu.Unlock()
	t.mu.Lock()
	t.count++
}

// SubCount 减少任务池的任务数量
func (t *TaskPool) SubCount() {
	defer t.mu.Unlock()
	t.mu.Lock()
	t.count--
}

// GetCount 获取任务池的任务数量
func (t *TaskPool) GetCount() int {
	defer t.mu.RUnlock()
	t.mu.RLock()
	return t.count
}

func (t *TaskPool) GetFails() int {
	defer t.mu.RUnlock()
	t.mu.RLock()
	return t.fails
}

func (t *TaskPool) AddFails() {
	defer t.mu.Unlock()
	t.mu.Lock()
	t.fails++
}

// Close 关闭任务池
// 请不要重复调用，否则会panic（也可以使用sync.Once，但有时候 我们不应该过度保姆式 设计）
func (t *TaskPool) Close() {
	for t.GetCount() != 0 {
		time.Sleep(time.Second * 1)
	}
	close(t.closeChan)
}
