FROM golang:latest

WORKDIR /driver_build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal ./internal
COPY cmd ./cmd

RUN go build -o app /driver_build/cmd/DriverService/main.go
