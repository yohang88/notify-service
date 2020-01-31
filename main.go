package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yohang88/notify-service/queue"
)

func main() {
	queue.Init("amqp://localhost")

	if os.Args[1] == "worker" {
		worker()
	} else {
		publisher()
	}
}

func publisher() {
	for {
		if err := queue.Publish("push_message", []byte(`{"num":6.13,"data":["a","b"]}`)); err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func worker() {
	// obtain the channel which we subscribe to
	msgs, close, err := queue.Subscribe("push_message")
	if err != nil {
		panic(err)
	}
	defer close()
	forever := make(chan bool)

	go func() {
		// Receive messages from the channel forever
		for d := range msgs {
			// then print the result to STDOUT, along with the time
			fmt.Println(time.Now().Format("01-02-2006 15:04:05"), byteToString(d.Body))
			// acknowledge the message so that it is cleared from the queue
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func byteToString(b []byte) string {
	s := string(b)

	return s
}
