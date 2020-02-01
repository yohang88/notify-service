package main

import (
    "encoding/json"
    "fmt"
    "github.com/sirupsen/logrus"
    "github.com/yohang88/notify-service/queue"
    "net/http"
    "net/url"
    "os"
)

type Message struct {
    Phone string
    Content string
}

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
            myJsonString := byteToString(d.Body)

            fmt.Println(myJsonString)

            var message Message
            json.Unmarshal([]byte(myJsonString), &message)

            sendSms(message.Phone, message.Content)

            logrus.WithFields(logrus.Fields{"event_name": "QUEUE_WORKER_WORKING", "number": message.Phone, "content": message.Content}).
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

func sendSms(phone string, content string) {
    baseUrl, _ := url.Parse(os.Getenv("SMS_BASE_URL"))

    params := url.Values{}
    params.Add("username", os.Getenv("SMS_USERNAME"))
    params.Add("key", os.Getenv("SMS_API_KEY"))
    params.Add("number", phone)
    params.Add("message", content)

    baseUrl.RawQuery = params.Encode()

    _, err := http.Get(baseUrl.String())

    if err != nil {
        panic(err)
    }
}