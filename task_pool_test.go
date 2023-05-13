package pool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestTaskPool(t *testing.T) {
	pool := NewTaskPool(
		WithPoolSize(30),
		WithBufferSize(10))

	t.Run("test task pool", func(t *testing.T) {
		// 创建200个任务
		for i := 0; i < 200; i++ {
			ctx, _ := context.WithTimeout(context.Background(), 10)
			err := pool.AddTask(Task, ctx)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

}

func Task() {
	fmt.Println("hello world")
	time.Sleep(time.Second * 10)
}
