package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
)

var async = flag.Bool("a", false, "use async")

func main() {
	flag.Parse()

	ctx := context.Background()

	logger := log.Default()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:     []string{"127.0.0.1:29092"},
		Topic:       "read_topic",
		Async:       *async,
		Logger:      kafka.LoggerFunc(logger.Printf),
		ErrorLogger: kafka.LoggerFunc(logger.Printf),
		BatchSize:   2000,
	})
	defer writer.Close()
	testMsgpref := "{\n    \"id\": \"284665d6-0190-49e7-34e9-9b4061acc261\",\n    \"source\": \"/trip\",\n    \"type\": \"trip.event.created\",\n    \"datacontenttype\": \"application/json\",\n    \"time\": \"2023-11-09T17:31:00Z\",\n    \"data\": {\n        \"trip_id\": \""
	testMsgpostf := "w83fe42d6-b85f-4e2a-93a2-858513acb148\",\n        \"offer_id\": \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN0cmluZyIsImZyb20iOnsibGF0IjowLCJsbmciOjB9LCJ0byI6eyJsYXQiOjAsImxuZyI6MH0sImNsaWVudF9pZCI6InN0cmluZyIsInByaWNlIjp7ImFtb3VudCI6OTkuOTUsImN1cnJlbmN5IjoiUlVCIn19.fg0Bv2ONjT4r8OgFqJ2tpv67ar7pUih2LhDRCRhWW3c\",\n        \"price\": {\n            \"currency\": \"RUB\",\n            \"amount\": 100\n        },\n        \"status\": \"DRIVER_SEARCH\",\n        \"from\": {\n            \"lat\": 0,\n            \"lng\": 0\n        },\n        \"to\": {\n            \"lat\": 0,\n            \"lng\": 0\n        }\n    }\n}"
	err := writer.WriteMessages(ctx, kafka.Message{Key: []byte(strconv.Itoa(0)), Value: []byte(testMsgpref + strconv.Itoa(0) + testMsgpostf)})
	//for i := 0; i < 524_288; i++ {
	//	var a int
	//	fmt.Scanf("%d", &a)
	//	err := writer.WriteMessages(ctx, kafka.Message{Key: []byte(strconv.Itoa(i)), Value: []byte(testMsgpref + strconv.Itoa(i) + testMsgpostf)})
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	if err != nil {
		fmt.Println(err.Error())
	}
}
