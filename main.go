package main

import (
    "strconv"
    "time"

    "github.com/sirupsen/logrus"
    "github.com/yohang88/notify-service/queue"
)

func main() {
    logrus.SetFormatter(&logrus.JSONFormatter{})

    queue.Init("amqp://localhost")

    publisher()
}

func publisher() {
    for {
        messagesCount, err := queue.Stats("push_message")
        if err != nil {
            panic(err)
        }

        logrus.WithFields(logrus.Fields{"event_name": "QUEUE_PUSH"}).
            Info("Message pushed to queues. Count: " + strconv.Itoa(messagesCount))

        if err := queue.Publish("push_message", []byte(`{"num":6.13,"data":["a","b"]}`)); err != nil {
            panic(err)
        }

        time.Sleep(500 * time.Millisecond)
    }
}
