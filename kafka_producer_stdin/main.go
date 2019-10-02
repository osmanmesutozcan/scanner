package main

import (
	"bufio"
	"flag"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"os"
)

func main() {
	flag.Parse()
	topic := flag.Arg(0)
	if len(topic) == 0 {
		panic("Cannot have empty topic")
	}

	sc := bufio.NewScanner(os.Stdin)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	for sc.Scan() {
		message := sc.Text()
		_ = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
		}, nil)
	}
}
