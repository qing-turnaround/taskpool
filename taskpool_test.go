package pool

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestTaskPool(t *testing.T) {

	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			var (
				mu    = new(sync.RWMutex)
				count = 0
			)

			// 定义任务函数
			Task := func() error {
				mu.Lock()
				count++
				mu.Unlock()
				return nil
			}

			// 创建一个任务池
			pool := NewTaskPool(
				WithPoolSize(40),
				WithBufferSize(10))
			var n = 10000
			// 创建200个任务
			for i := 0; i < n; i++ {
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
			if n != count+pool.GetFails() {
				t.Errorf("任务池有误")
			}
		})
	}

}
