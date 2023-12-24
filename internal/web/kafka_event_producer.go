package web

import (
	"context"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/svc"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type KafkaEventProducer struct {
	Writer  KafkaWriter
	timeout time.Duration
}

type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

const (
	KAFKA_MESSAGE_TYPE_PREFIX = "trip.command."
	DATA_CONTENT_TYPE         = "application/json"
)

func NewKafkaEventProducer(cfg *config.KafkaConfig) *KafkaEventProducer {
	logger := log.Default()
	return &KafkaEventProducer{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:     []string{cfg.Host + ":" + cfg.Port},
			Topic:       cfg.WriteTopic,
			Async:       cfg.AsyncWrite,
			Logger:      kafka.LoggerFunc(logger.Printf),
			ErrorLogger: kafka.LoggerFunc(logger.Printf),
			BatchSize:   cfg.WriteBatchSize,
		}),
		timeout: time.Duration(cfg.TimeoutSec) * time.Second,
	}
}

func (p *KafkaEventProducer) SendTripEvent(baseCtx context.Context, msg *svc.TripMessagePayload) {
	ctx, span := svc.Tracer.Start(baseCtx, "/send-trip-event")
	defer span.End()
	logger.Main.Info(fmt.Sprintf("Start sending trip event message\n%s", msg.ToString()))
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	message, err := json.Marshal(NewTripMessage(msg))
	if err != nil {
		logger.Main.Error(err.Error())
		return
	}
	err = p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.Id),
		Value: message,
	})
	if err != nil {
		logger.Main.Error(err.Error())
	}
	logger.Main.Info(fmt.Sprintf("Trip event message with id: %s successfully sended\n%s", msg.Id, msg.ToString()))
}

func (p *KafkaEventProducer) Free() {
	p.Writer.Close()
}

func NewTripMessage(msg *svc.TripMessagePayload) *TripMessage {
	return &TripMessage{
		Id:              msg.Id,
		Source:          msg.Source,
		Type:            KAFKA_MESSAGE_TYPE_PREFIX + string(msg.Status),
		DataContentType: DATA_CONTENT_TYPE,
		Time:            msg.Time.Format(time.RFC3339),
		Data: TripData{
			TripId: msg.TripId,
		},
	}
}

type TripMessage struct {
	Id              string   `json:"id"`
	Source          string   `json:"source"`
	Type            string   `json:"type"`
	DataContentType string   `json:"datacontenttype"`
	Time            string   `json:"time"`
	Data            TripData `json:"data"`
}

type TripData struct {
	TripId string `json:"trip_id"`
}
