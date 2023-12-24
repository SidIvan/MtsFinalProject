package web

import (
	"context"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/test/mock"
	"driver-service/internal/web"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/segmentio/kafka-go"
	"golang.org/x/tools/container/intsets"
	"testing"
	"time"
)

const (
	KAFKA_PORT = "9092/tcp"
)

var (
	kafkaEventConsumer *web.CreateEventConsumer
	kafkaReader        *mock.MockKafkaReader
	testMessage        = kafka.Message{
		Topic:         "test_topic",
		Partition:     1,
		Offset:        2,
		HighWaterMark: 3,
		Key:           nil,
		Value:         []byte("{\n    \"id\": \"284655d6-0190-49e7-34e9-9b4060acc261\",\n    \"source\": \"/trip\",\n    \"type\": \"trip.event.created\",\n    \"datacontenttype\": \"application/json\",\n    \"time\": \"2023-11-09T17:31:00Z\",\n    \"data\": {\n        \"trip_id\": \"e82c42d6-b86f-4e2a-93a2-858413acb148\",\n        \"offer_id\": \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN0cmluZyIsImZyb20iOnsibGF0IjowLCJsbmciOjB9LCJ0byI6eyJsYXQiOjAsImxuZyI6MH0sImNsaWVudF9pZCI6InN0cmluZyIsInByaWNlIjp7ImFtb3VudCI6OTkuOTUsImN1cnJlbmN5IjoiUlVCIn19.fg0Bv2ONjT4r8OgFqJ2tpv67ar7pUih2LhDRCRhWW3c\",\n        \"price\": {\n            \"currency\": \"RUB\",\n            \"amount\": 100\n        },\n        \"status\": \"DRIVER_SEARCH\",\n        \"from\": {\n            \"lat\": 0,\n            \"lng\": 0\n        },\n        \"to\": {\n            \"lat\": 0,\n            \"lng\": 0\n        }\n    }\n}"),
		Headers:       nil,
		WriterData:    nil,
		Time:          time.Now(),
	}
)

func initKafkaConsumerMocks(t *testing.T) {
	logger.InitTestLogger(t)
	kafkaConfig = &config.KafkaConfig{
		Host:              "",
		Port:              "",
		ReadTopic:         "test_topic",
		WriteTopic:        "",
		GroupId:           "",
		SessionTimeoutSec: 0,
		AsyncWrite:        false,
		WriteBatchSize:    0,
		TimeoutSec:        intsets.MaxInt,
	}
	ctrl := gomock.NewController(t)
	driverService = mock.NewMockDriverService(ctrl)
	driverService.EXPECT().CloseEventWriter()
	kafkaEventConsumer = web.NewCreateEventConsumer(kafkaConfig, driverService)
	kafkaEventConsumer.Stop()
	kafkaReader = mock.NewMockKafkaReader(ctrl)
	kafkaEventConsumer.Reader = kafkaReader
}

func TestReadAndServeMessageSuccess(t *testing.T) {
	initKafkaConsumerMocks(t)
	kafkaReader.EXPECT().ReadMessage(gomock.Any()).Return(testMessage, nil)
	driverService.EXPECT().CreateTrip(gomock.Any(), gomock.Any())
	kafkaEventConsumer.ReadAndServeMessage(context.Background())
}

func TestReadAndServeMessageFailed(t *testing.T) {
	initKafkaConsumerMocks(t)
	kafkaReader.EXPECT().ReadMessage(gomock.Any()).Return(testMessage, errors.New("test_error"))
	kafkaEventConsumer.ReadAndServeMessage(context.Background())
}
