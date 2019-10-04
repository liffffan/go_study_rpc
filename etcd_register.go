package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

func main() {
	// 连接 etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 获取租约，有效期5秒
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}

	// 设置 key 的时候把租约带上，这样的话续期就是5秒一次
	_, err = cli.Put(context.TODO(), "foo", "bar", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}

	// 永久续约
	ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}

	for {
		ka := <-ch
		fmt.Println("ttl:", ka)
	}

}
