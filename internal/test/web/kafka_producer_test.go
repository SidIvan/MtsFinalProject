package web

import (
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/svc"
	"driver-service/internal/test/mock"
	"driver-service/internal/web"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/segmentio/kafka-go"
	"golang.org/x/tools/container/intsets"
	"testing"
	"time"
)

const ()

var (
	kafkaWriter        *mock.MockKafkaWriter
	kafkaEventProducer *web.KafkaEventProducer
	kafkaConfig        = &config.KafkaConfig{
		Host:              "",
		Port:              "",
		ReadTopic:         "",
		WriteTopic:        "",
		GroupId:           "",
		SessionTimeoutSec: 0,
		AsyncWrite:        false,
		WriteBatchSize:    0,
		TimeoutSec:        intsets.MaxInt,
	}
	testMsgPayload = &svc.TripMessagePayload{
		Id:     "test_id",
		Source: "test_source",
		Time:   time.Now(),
		Status: svc.ACCEPTED,
		TripId: "test_trip_id",
	}
)

func initKafkaProducerMocks(t *testing.T) {
	logger.InitTestLogger(t)
	ctrl := gomock.NewController(t)
	kafkaWriter = mock.NewMockKafkaWriter(ctrl)
	kafkaEventProducer = web.NewKafkaEventProducer(kafkaConfig)
	kafkaEventProducer.Free()
	kafkaEventProducer.Writer = kafkaWriter
}

func TestSendTripEventSuccess(t *testing.T) {
	bytePayload, err := json.Marshal(web.NewTripMessage(testMsgPayload))
	if err != nil {
		t.Error(err.Error())
		return
	}
	expectedKafkaMessage := kafka.Message{
		Key:   []byte(testMsgPayload.Id),
		Value: bytePayload,
	}
	initKafkaProducerMocks(t)
	kafkaWriter.EXPECT().WriteMessages(gomock.Any(), expectedKafkaMessage).Return(nil)
	kafkaEventProducer.SendTripEvent(TEST_BASE_CONTEXT, testMsgPayload)
}

func TestSendTripFailed(t *testing.T) {
	bytePayload, err := json.Marshal(web.NewTripMessage(testMsgPayload))
	if err != nil {
		t.Error(err.Error())
		return
	}
	expectedKafkaMessage := kafka.Message{
		Key:   []byte(testMsgPayload.Id),
		Value: bytePayload,
	}
	initKafkaProducerMocks(t)
	kafkaWriter.EXPECT().WriteMessages(gomock.Any(), expectedKafkaMessage).Return(errors.New("test_error"))
	kafkaEventProducer.SendTripEvent(TEST_BASE_CONTEXT, testMsgPayload)
}
