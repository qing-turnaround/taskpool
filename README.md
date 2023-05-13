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
	for i := 0; i < 300; i++ {
		// 任务上下文, 10秒超时
		ctx, _ := context.WithTimeout(context.Background(), 20)
		// 添加任务
		err := pool.AddTask(Task, ctx)
		// 错误处理
		if err != nil {
			fmt.Println(err)
		}
	}

	// 等待10秒
	time.Sleep(time.Second * 10)
	// 关闭任务池
	pool.Close()
	fmt.Println(count)

}


// 定义任务函数
func Task() {
	mu.Lock()
	count++
	mu.Unlock()

}
```