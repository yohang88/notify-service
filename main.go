package main

import (
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/labstack/echo"

    "github.com/sirupsen/logrus"
    "github.com/yohang88/notify-service/queue"
)

func main() {
    logrus.SetFormatter(&logrus.JSONFormatter{})

    queue.Init("amqp://localhost")
    // publisher()

    e := echo.New()
    // e.Use(middleware.Logger())

    e.GET("/", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{"version": "1.0.0"})
    })

    e.POST("/broadcasts", createBroadcast)

    e.Logger.Fatal(e.Start(":8000"))
}

func createBroadcast(c echo.Context) error {
    err := queue.Publish("push_message", []byte(`{"num":6.13,"data":["a","b"]}`))

    if err != nil {
        log.Fatal(err)
    }

    logrus.WithFields(logrus.Fields{"event_name": "BROADCAST_CREATE"}).Info("Create broadcast")

    return c.NoContent(http.StatusNoContent)
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
