package main

import (
	"flag"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	flag.Parse()
	topic := flag.Arg(0)
	if len(topic) == 0 {
		panic("Cannot have empty topic")
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id": "stdoutgroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("%s\n", string(msg.Value))
		}
	}
}
