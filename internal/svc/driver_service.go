package svc

import (
	"context"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/model"
	"encoding/json"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.opentelemetry.io/otel"
	"sync"
	"time"
)

var Tracer = otel.Tracer("driver")

type DriverRepo interface {
	CreateTrip(context.Context, *model.Trip) string
	PutTripStatus(context.Context, string, string, model.TripStatus) bool
	GetTrip(context.Context, string, string) *model.Trip
	GetTrips(context.Context, string) []model.Trip
	SaveCancelReason(context.Context, string, string, string)
}

type LocationClient interface {
	GetDrivers(context.Context, *GetDriversPayload) []model.Driver
}

type EventProducer interface {
	SendTripEvent(context.Context, *TripMessagePayload)
	Free()
}

type DriverService interface {
	GetTrips(context.Context, string) []model.Trip
	GetTrip(context.Context, string, string) *model.Trip
	CancelTrip(context.Context, string, string, string) ChangeTripStatusStatus
	AcceptTrip(context.Context, string, string) ChangeTripStatusStatus
	StartTrip(context.Context, string, string) ChangeTripStatusStatus
	EndTrip(context.Context, string, string) ChangeTripStatusStatus
	CreateTrip(context.Context, model.Trip)
	CloseEventWriter()
}

type Notificationer interface {
	SendNotificationToDriver(model.Driver)
}

type DriverServiceImpl struct {
	source         string
	radius         float64
	driverRepo     DriverRepo
	locationClient LocationClient
	eventProducer  EventProducer
	notificationer Notificationer
}

func NewDriverService(cfg *config.DriverServiceConfig, driverRepo DriverRepo, locationClient LocationClient, eventProducer EventProducer) *DriverServiceImpl {
	return &DriverServiceImpl{
		source:         cfg.Source,
		radius:         cfg.SearchRadius,
		driverRepo:     driverRepo,
		locationClient: locationClient,
		eventProducer:  eventProducer,
	}
}

func (s DriverServiceImpl) GetTrips(ctx context.Context, userId string) []model.Trip {
	ctx, span := Tracer.Start(ctx, "/get-trips-svc-method")
	defer span.End()
	return s.driverRepo.GetTrips(ctx, userId)
}

func (s DriverServiceImpl) GetTrip(ctx context.Context, userId string, tripId string) *model.Trip {
	ctx, span := Tracer.Start(ctx, "/get-trip-svc-method")
	defer span.End()
	return s.driverRepo.GetTrip(ctx, tripId, userId)
}

type ChangeTripStatusStatus byte

const (
	Ok               = ChangeTripStatusStatus(iota)
	TripNotFound     = ChangeTripStatusStatus(iota)
	InvalidNewStatus = ChangeTripStatusStatus(iota)
	InternalError    = ChangeTripStatusStatus(iota)
)

func (s DriverServiceImpl) CancelTrip(ctx context.Context, userId string, tripId string, reason string) ChangeTripStatusStatus {
	ctx, span := Tracer.Start(ctx, "/cancel-trip-svc-method")
	defer span.End()
	trip := s.driverRepo.GetTrip(ctx, tripId, userId)
	if trip == nil {
		logger.Main.Warn(fmt.Sprintf("Trip with id: %s not found for driver with id: %s", tripId, userId))
		return TripNotFound
	}
	if model.IsValidChangeStatus(trip.Status, model.CANCELED) {
		if !s.driverRepo.PutTripStatus(ctx, tripId, userId, model.CANCELED) {
			logger.Main.Error(fmt.Sprintf("Correct cancel trip operation failed"))
			return InternalError
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			s.driverRepo.SaveCancelReason(ctx, tripId, userId, reason)
			wg.Done()
		}()
		wg.Wait()
		return Ok
	}
	logger.Main.Warn(fmt.Sprintf("Invalid attempt to change trip status from %s to %s", trip.Status, model.CANCELED))
	return InvalidNewStatus
}

func (s DriverServiceImpl) AcceptTrip(ctx context.Context, userId string, tripId string) ChangeTripStatusStatus {
	ctx, span := Tracer.Start(ctx, "/accept-trip-svc-method")
	defer span.End()
	trip := s.driverRepo.GetTrip(ctx, tripId, userId)
	if trip == nil {
		logger.Main.Warn(fmt.Sprintf("Trip with id: %s not found for driver with id: %s", tripId, userId))
		return TripNotFound
	}
	if !model.IsValidChangeStatus(trip.Status, model.DRIVER_FOUND) {
		logger.Main.Warn(fmt.Sprintf("Invalid attempt to change trip status from %s to %s", trip.Status, model.DRIVER_FOUND))
		return InvalidNewStatus
	}
	if s.driverRepo.PutTripStatus(ctx, tripId, userId, model.DRIVER_FOUND) {
		go s.sendAcceptedTripMessage(context.Background(), NewTripMessagePayload(s.genMsgId(), s.source, tripId))
		return Ok
	}
	logger.Main.Error(fmt.Sprintf("Correct accept trip operation failed"))
	return InternalError
}

func (s DriverServiceImpl) StartTrip(ctx context.Context, userId string, tripId string) ChangeTripStatusStatus {
	ctx, span := Tracer.Start(ctx, "/start-trip-svc-method")
	defer span.End()
	trip := s.driverRepo.GetTrip(ctx, tripId, userId)
	if trip == nil {
		return TripNotFound
	}
	if !model.IsValidChangeStatus(trip.Status, model.STARTED) {
		logger.Main.Warn(fmt.Sprintf("Invalid attempt to change trip status from %s to %s", trip.Status, model.STARTED))
		return InvalidNewStatus
	}
	if s.driverRepo.PutTripStatus(ctx, tripId, userId, model.STARTED) {
		go s.sendStartedTripMessage(context.Background(), NewTripMessagePayload(s.genMsgId(), s.source, tripId))
		return Ok
	}
	logger.Main.Error(fmt.Sprintf("Correct start trip operation failed"))
	return InternalError
}

func (s DriverServiceImpl) EndTrip(ctx context.Context, userId string, tripId string) ChangeTripStatusStatus {
	ctx, span := Tracer.Start(ctx, "/end-trip-svc-method")
	defer span.End()
	trip := s.driverRepo.GetTrip(ctx, tripId, userId)
	if trip == nil {
		return TripNotFound
	}
	if !model.IsValidChangeStatus(trip.Status, model.ENDED) {
		logger.Main.Warn(fmt.Sprintf("Invalid attempt to change trip status from %s to %s", trip.Status, model.ENDED))
		return InvalidNewStatus
	}
	if s.driverRepo.PutTripStatus(ctx, tripId, userId, model.ENDED) {
		go s.sendEndedTripMessage(context.Background(), NewTripMessagePayload(s.genMsgId(), s.source, tripId))
		return Ok
	}
	logger.Main.Error(fmt.Sprintf("Correct end trip operation failed"))
	return InternalError
}

func (s DriverServiceImpl) CreateTrip(ctx context.Context, trip model.Trip) {
	ctx, span := Tracer.Start(ctx, "/create-trip-svc-method")
	defer span.End()
	trip.Status = model.CREATED
	insertedId := s.driverRepo.CreateTrip(ctx, &trip)
	if insertedId == "" {
		logger.Main.Warn("Trip does not created")
	} else {
		logger.Main.Debug(fmt.Sprintf("Created trip with id: %s", insertedId))
	}
	drivers := s.locationClient.GetDrivers(ctx, &GetDriversPayload{
		LatLngLiteral: trip.From,
		Radius:        s.radius,
	})
	for i := 0; i < len(drivers); i++ {
		s.notificationer.SendNotificationToDriver(drivers[i])
	}
}

type MessageStatus string

const (
	ACCEPTED = MessageStatus("accepted")
	STARTED  = MessageStatus("started")
	ENDED    = MessageStatus("ended")
)

func (s DriverServiceImpl) sendAcceptedTripMessage(ctx context.Context, msgPayload *TripMessagePayload) {
	ctx, span := Tracer.Start(ctx, "/send-accepted-trip-message-method")
	defer span.End()
	s.sendTripMessage(ctx, msgPayload, ACCEPTED)
}

func (s DriverServiceImpl) sendStartedTripMessage(ctx context.Context, msgPayload *TripMessagePayload) {
	ctx, span := Tracer.Start(ctx, "/send-started-trip-message-method")
	defer span.End()
	s.sendTripMessage(ctx, msgPayload, STARTED)
}

func (s DriverServiceImpl) sendEndedTripMessage(ctx context.Context, msgPayload *TripMessagePayload) {
	ctx, span := Tracer.Start(ctx, "/send-ended-trip-message-method")
	defer span.End()
	s.sendTripMessage(ctx, msgPayload, ENDED)
}

func (s DriverServiceImpl) sendTripMessage(ctx context.Context, msgPayload *TripMessagePayload, status MessageStatus) {
	if msgPayload.Id == "" {
		return
	}
	msgPayload.Status = status
	s.eventProducer.SendTripEvent(ctx, msgPayload)
}

func (s DriverServiceImpl) genMsgId() string {
	var err error
	for i := 0; i < 10; i++ {
		id, err := gonanoid.New()
		if err != nil {
			continue
		}
		return id
	}
	logger.Main.Error(err.Error())
	return ""
}

func (s DriverServiceImpl) CloseEventWriter() {
	s.eventProducer.Free()
}

type TripMessagePayload struct {
	Id     string        `json:"id"`
	Source string        `json:"source"`
	Time   time.Time     `json:"time"`
	Status MessageStatus `json:"status"`
	TripId string        `json:"trip_id"`
}

func NewTripMessagePayload(id string, source string, tripId string) *TripMessagePayload {
	return &TripMessagePayload{
		Id:     id,
		Source: source,
		Time:   time.Now(),
		TripId: tripId,
	}
}

func (p *TripMessagePayload) ToString() string {
	jsonData, err := json.Marshal(p)
	if err != nil {
		logger.Main.Error(err.Error())
	}
	return string(jsonData)
}
