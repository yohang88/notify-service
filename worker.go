package main

import (
    "github.com/sirupsen/logrus"
    "github.com/yohang88/notify-service/queue"
)

func main() {
    logrus.SetFormatter(&logrus.JSONFormatter{})

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
            logrus.WithFields(logrus.Fields{"event_name": "QUEUE_WORKER_WORKING", "event_data": byteToString(d.Body)}).
                Info("New message.")

            d.Ack(false)
        }
    }()

    logrus.WithFields(logrus.Fields{"event_name": "QUEUE_WORKER_READY"}).
        Info("Waiting for messages")

    <-forever
}

func byteToString(b []byte) string {
    s := string(b)

    return s
}
