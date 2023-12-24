package web

import (
	"context"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/model"
	"driver-service/internal/svc"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type CreateEventConsumer struct {
	driverService svc.DriverService
	Reader        KafkaReader
}

type KafkaReader interface {
	ReadMessage(context.Context) (kafka.Message, error)
	Close() error
}

func NewCreateEventConsumer(cfg *config.KafkaConfig, driverService svc.DriverService) *CreateEventConsumer {
	return &CreateEventConsumer{
		driverService: driverService,
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{cfg.Host + ":" + cfg.Port},
			Topic:          cfg.ReadTopic,
			GroupID:        cfg.GroupId,
			SessionTimeout: time.Duration(cfg.SessionTimeoutSec) * time.Second,
		}),
	}
}

func (c *CreateEventConsumer) ReadAndServeMessage(ctx context.Context) {
	ctx, span := svc.Tracer.Start(ctx, "/read-and-serve-message-action")
	defer span.End()
	msg, err := c.Reader.ReadMessage(ctx)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var msgData CreateTripMessage
	err = json.Unmarshal(msg.Value, &msgData)
	if err != nil {
		log.Println(err.Error())
		return
	}
	logger.Main.Debug(fmt.Sprintf("Read create trip event\n%s", string(msg.Value)))
	trip := model.Trip{
		Id:       msgData.Data.TripId,
		DriverId: "",
		From:     msgData.Data.From,
		To:       msgData.Data.To,
		Price:    msgData.Data.Price,
		Status:   "",
	}
	c.driverService.CreateTrip(ctx, trip)
}

func (c *CreateEventConsumer) Start(ctx context.Context) {
	for {
		c.ReadAndServeMessage(ctx)
	}
}

func (c *CreateEventConsumer) Stop() {
	c.Reader.Close()
	c.driverService.CloseEventWriter()
}

type CreateTripMessage struct {
	Id              string         `json:"id"`
	Source          string         `json:"source"`
	Type            string         `json:"type"`
	DataContentType string         `json:"datacontenttype"`
	Time            string         `json:"time"`
	Data            CreateTripData `json:"data"`
}

type CreateTripData struct {
	TripId  string              `json:"trip_id"`
	OfferId string              `json:"offer_id"`
	Price   model.Money         `json:"price"`
	Status  model.TripStatus    `json:"status"`
	From    model.LatLngLiteral `json:"from"`
	To      model.LatLngLiteral `json:"to"`
}
