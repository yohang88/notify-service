package main

import (
	"fmt"
	"time"

	"github.com/yohang88/notify-service/queue"
)

func main() {
	queue.Init("amqp://localhost")

	worker()
}

func worker() {
	messages, close, err := queue.Subscribe("push_message")

	if err != nil {
		panic(err)
	}

	defer close()
	forever := make(chan bool)

	go func() {
		for d := range messages {
			fmt.Println(time.Now().Format("2006-01-02T15:04:05-0700"), byteToString(d.Body))
			d.Ack(false)
		}
	}()

	fmt.Println(time.Now().Format("2006-01-02T15:04:05-0700"), "Waiting for messages")
	<-forever
}

func byteToString(b []byte) string {
	s := string(b)

	return s
}
