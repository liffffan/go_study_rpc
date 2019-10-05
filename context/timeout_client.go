package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	TranceId = "trace_id"
)

var (
	wg sync.WaitGroup
)

type ResPack struct {
	r   *http.Response
	err error
}

func work(ctx context.Context) {
	// Transport : HTTP 底层通讯对象
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	defer wg.Done()
	// 创建一个管道，类型是服务端返回的结果
	c := make(chan ResPack, 1)
	// 启用一个 HTTP 的 Request
	req, _ := http.NewRequest("GET", "http://localhost:9200", nil)
	// 通过 goroutine 来发送 http 请求
	go func() {
		// 如果请求有结果返回就放入管道里
		resp, err := client.Do(req)
		pack := ResPack{r: resp, err: err}
		c <- pack
	}()

	// 超时处理
	select {
	case <-ctx.Done():
		// 如果超时了，会返回错误。取消请求，把错误信号直接丢掉
		tr.CancelRequest(req)
		<-c
		fmt.Println("Timeout!")
	case res := <-c:
		// 如果没有超时，就会走这个分支
		if res.err != nil {
			fmt.Println(res.err)
			return
		}
		a(ctx)
		// 这个不关会导致内存泄露：https://segmentfault.com/a/1190000020086816?utm_source=tag-newest
		defer res.r.Body.Close()
		out, _ := ioutil.ReadAll(res.r.Body)
		fmt.Printf("Server Response:%s", out)
	}
}

// 传递参数
func a(ctx context.Context) {
	trance_id := ctx.Value(TranceId)
	fmt.Printf("Trace_id is: %v, process of a\n", trance_id)
	b(ctx)
}

func b(ctx context.Context) {
	trance_id := ctx.Value(TranceId)
	fmt.Printf("Trace_id is: %v, process of b\n", trance_id)
	c(ctx)
}

func c(ctx context.Context) {
	trance_id := ctx.Value(TranceId)
	fmt.Printf("Trace_id is: %v, process of c\n", trance_id)
}

func main() {
	ctx := context.WithValue(context.Background(), TranceId, rand.Int63())
	a(ctx)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	wg.Add(1)
	go work(ctx)
	wg.Wait()
	fmt.Println("Finished")
}
