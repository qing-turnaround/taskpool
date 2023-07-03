# taskPool
* 用Go语言实现一个goroutine pool（任务池）

# 如何安装
* `go get github.com/xing-you-ji/taskpool`


# Example
```go
package main

import (
	"context"
	"fmt"
	pool "github.com/xing-you-ji/taskpool"
	"sync"
	"time"
)

var (
	mu  = new(sync.RWMutex)
	count  = 0
)


func main() {
	// 创建一个任务池
	pool := pool.NewTaskPool(
		pool.WithPoolSize(40),
		pool.WithBufferSize(10))

	// 创建200个任务
	for i := 0; i < 10000; i++ {
		// 任务上下文, 10秒超时
		ctx, _ := context.WithTimeout(context.Background(), 20)
		// 添加任务
		err := pool.AddTask(Task, ctx)
		// 错误处理
		if err != nil {
			fmt.Println(err)
		}
	}

	// 关闭任务池
	pool.Close()
	fmt.Printf("任务池执行完毕，失败任务数量：%d\n", pool.GetFails())
	fmt.Printf("任务池执行完毕，结果为：%d\n", count)
}


// 定义任务函数
func Task() error {
	mu.Lock()
	count++
	mu.Unlock()
	return nil
}
```