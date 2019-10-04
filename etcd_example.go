package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	testConnect()
	testPut()
	testWatch()
}

func testConnect() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect success")
	defer cli.Close()
}

func testPut() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect success")
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "/logagent/conf/", "sample_value")
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	fmt.Println("put success")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/logagent/conf/")
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("get %s, %s\n", ev.Key, ev.Value)
	}

}

func testWatch() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect success")
	defer cli.Close()

	// 开一个 goroutine 进行 etcd 配置修改，测试 watch 机制
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err = cli.Put(ctx, "/logagent/conf/", "test_watch")
		cancel()
		if err != nil {
			fmt.Println("put failed, err:", err)
			return
		}
	}()

	for {
		rch := cli.Watch(context.Background(), "/logagent/conf/")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("watch success, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}

}
