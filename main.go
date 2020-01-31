package main

import (
	"time"

	"github.com/yohang88/notify-service/queue"
)

func main() {
	queue.Init("amqp://localhost")

	publisher()
}

func publisher() {
	for {
		if err := queue.Publish("push_message", []byte(`{"num":6.13,"data":["a","b"]}`)); err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}