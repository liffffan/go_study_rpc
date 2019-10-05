package main

import (
	"context"
	"log"
	"os"
	"time"
)

var logg2 *log.Logger

func timeoutHandler() {
	//ctx, cancel := context.WithTimeout(context.Background(),  2*time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	go doStuff2(ctx)
	time.Sleep(10 * time.Second)
	cancel()
}

// 每一秒 work 一下，同时会判断 ctx 是否被取消了，如果是就退出
// 也就是在 cancel 执行之前，这个 ctx.Done() 永远都不会触发
// cancel 在被调用的时候会往 ctx.Done() 这条管道插入一条消息
func doStuff2(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			logg2.Printf("done")
			return
		default:
			logg2.Printf("work")
		}
	}
}

func main() {
	// 初始化一个系统日志
	logg2 = log.New(os.Stdout, "", log.Ltime)
	timeoutHandler()
	logg2.Printf("down")
}
