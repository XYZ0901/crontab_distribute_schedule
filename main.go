package main

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func cr() {
	c := cron.New(cron.WithSeconds())
	id, err := c.AddFunc("*/20 * * * * *", func() {
		fmt.Println("hello")
	})
	if err != nil {
		log.Fatalln(err)
	}
	c.Start()

	e := c.Entry(id)

	<-time.NewTimer(e.Next.Sub(time.Now())).C
	fmt.Println(time.Now())
	c.AddFunc("*/5 * * * * *", func() {
		fmt.Println("hello2")
	})
}

func main() {
	ch := make(chan struct{})

	//cancelMap := make(map[string]chan struct{})
	// 任务启动 -> 往cancelMap注册id switch case?
	// kill任务 -> 通知chan
	//clientv3.NewKV()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("ctx.done")
		default:
			for i := 0; i < 10; i++ {
				fmt.Println(i)
				time.Sleep(time.Second)
			}
		}
	}()

	time.Sleep(3 * time.Second)
	cancel()

	<-ch
}
