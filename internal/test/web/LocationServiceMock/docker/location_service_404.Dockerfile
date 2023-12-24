FROM golang:latest

WORKDIR /mock

COPY . .

RUN go mod download
RUN go build -o mockSvc

EXPOSE 8081

ENTRYPOINT ./mockSvc -mode 404