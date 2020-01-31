# Simple Messaging Service

## Run RabbitMQ
```
docker run -d --hostname my-rabbit --name my-rabbit -p 5672:5672 rabbitmq:3
```

## Run Queue/Job Worker
```
$ go build worker.go
$ ./worker

2020-02-01T00:35:47+0700 Waiting for messages
```

## Run Publisher
```
$ go build main.go
$ ./main
```
