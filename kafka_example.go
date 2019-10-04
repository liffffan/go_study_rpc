package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	"sync"
	"time"
)

func main() {
	Producer()

}

func Producer() {
	config := sarama.NewConfig()
	// Ack 发送消息后需不需要回应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 分片策略
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	msg := &sarama.ProducerMessage{}
	msg.Topic = "nginx_log"
	msg.Value = sarama.StringEncoder("this is a good test, my message is good")

	client, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	defer client.Close()

	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed, ", err)
		return
	}

	fmt.Printf("pid:%v offset:%v \n", pid, offset)
}

// 消费者代码
var (
	wg sync.WaitGroup
)

func customer() {
	consumer, err := sarama.NewConsumer(strings.Split("localhost:9092", ","), nil)
	if err != nil {
		fmt.Println("Failed to start consumer:%s", err)
		return
	}

	// 指定消费的 Topic 名称，因为生产会分片，所以消费也是按分片来消费，这里是拿到这个 Topic 的所有分片列表
	partitionList, err := consumer.Partitions("nginx_log")
	if err != nil {
		fmt.Println("Failed to get the list of partitions :%s", err)
		return
	}
	fmt.Println(partitionList)

	// 每个分片起一个 goroutine 来消费
	for partition := range partitionList {
		// 传入 Topic 名称、分片序号、从什么位置进行消费
		pc, err := consumer.ConsumePartition("nginx_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("Failed to start consumer for partition %d:%s\n", partition, err)
			return
		}

		defer pc.AsyncClose()
		wg.Add(1)
		go func(pc1 sarama.PartitionConsumer) {
			for msg := range pc1.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
			wg.Done()
		}(pc)
	}

	wg.Wait()

	//time.Sleep(time.Hour)
	consumer.Close()
}
