package main

import (
    "archive/zip"
    "log"
    "mime/multipart"
    "net/http"
    "strconv"
    "time"

    "github.com/labstack/echo"
    "github.com/sirupsen/logrus"
    "github.com/tealeg/xlsx"
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
    content := c.FormValue("content")
    file, err := c.FormFile("file")

    if err != nil {
        return err
    }

    records := readAndParseXls(file)

    for _, record := range records {
        // Push to Queues
        err = queue.Publish("push_message", []byte(`{"phone":"`+record+`","content":"`+content+`"}`))

        if err != nil {
            log.Fatal(err)
        }
    }

    logrus.WithFields(logrus.Fields{"event_name": "BROADCAST_CREATE", "content": content}).Info("Create broadcast")

    return c.NoContent(http.StatusNoContent)
}

func readAndParseXls(file *multipart.FileHeader) []string {
    src, err := file.Open()
    if err != nil {
        log.Fatal(err)
    }

    defer src.Close()

    reader, err := zip.NewReader(src, file.Size)
    xlFile, err := xlsx.ReadZipReader(reader)

    if err != nil {
        log.Fatal(err)
    }

    var records []string

    for _, sheet := range xlFile.Sheets {
        for _, row := range sheet.Rows {
            for _, cell := range row.Cells {
                text := cell.String()

                records = append(records, text)
            }
        }
    }

    return records
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
