# build stage
FROM golang:alpine AS build-env
WORKDIR /app
RUN apk --no-cache add build-base git gcc
ADD . /app
RUN go get -d . && go build worker.go

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/worker /app/
ENTRYPOINT ./worker
