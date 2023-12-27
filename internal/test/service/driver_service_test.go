package service

import (
	"context"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/model"
	"driver-service/internal/svc"
	"driver-service/internal/test/mock"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

var (
	testDriverServiceConfig = config.DriverServiceConfig{
		Source: "test_source",
	}
	testTrips = []model.Trip{
		{
			Id:       "id_1",
			DriverId: "driver_id_1",
			From: model.LatLngLiteral{
				Lat: 1,
				Lng: 1,
			},
			To: model.LatLngLiteral{
				Lat: 1.5,
				Lng: 1.5,
			},
			Price: model.Money{
				Amount:   1,
				Currency: "CAD",
			},
			Status: model.CREATED,
		},
		{
			Id:       "id_2",
			DriverId: "driver_id_2",
			From: model.LatLngLiteral{
				Lat: 2,
				Lng: 2,
			},
			To: model.LatLngLiteral{
				Lat: 2.5,
				Lng: 2.5,
			},
			Price: model.Money{
				Amount:   2,
				Currency: "RUB",
			},
			Status: model.DRIVER_FOUND,
		},
		{
			Id:       "id_3",
			DriverId: "driver_id_3",
			From: model.LatLngLiteral{
				Lat: 3,
				Lng: 3,
			},
			To: model.LatLngLiteral{
				Lat: 3.5,
				Lng: 3.5,
			},
			Price: model.Money{
				Amount:   3,
				Currency: "USD",
			},
			Status: model.ON_POSITION,
		},
		{
			Id:       "id_4",
			DriverId: "driver_id_4",
			From: model.LatLngLiteral{
				Lat: 4,
				Lng: 4,
			},
			To: model.LatLngLiteral{
				Lat: 4.5,
				Lng: 4.5,
			},
			Price: model.Money{
				Amount:   4,
				Currency: "ALL",
			},
			Status: model.STARTED,
		},
		{
			Id:       "id_5",
			DriverId: "driver_id_5",
			From: model.LatLngLiteral{
				Lat: 5,
				Lng: 5,
			},
			To: model.LatLngLiteral{
				Lat: 5.5,
				Lng: 5.5,
			},
			Price: model.Money{
				Amount:   5,
				Currency: "TRY",
			},
			Status: model.ENDED,
		},
		{
			Id:       "id_6",
			DriverId: "driver_id_6",
			From: model.LatLngLiteral{
				Lat: 6,
				Lng: 6,
			},
			To: model.LatLngLiteral{
				Lat: 6.5,
				Lng: 6.5,
			},
			Price: model.Money{
				Amount:   6,
				Currency: "MOP",
			},
			Status: model.CANCELED,
		},
	}
)

const (
	TEST_DRIVER_ID            = "test_driver_id"
	TEST_EXISTING_TRIP_ID     = "test_trip_id"
	TEST_NON_EXISTING_TRIP_ID = "test_non_existing_trip_id"
	TEST_CANCEL_REASON        = "test_cancel_reason"
)

func initMocksAndDriverService(t *testing.T) (*svc.DriverServiceImpl, *mock.MockDriverRepo, *mock.MockLocationClient, *mock.MockEventProducer) {
	logger.InitTestLogger(t)
	ctrl := gomock.NewController(t)
	driverRepo := mock.NewMockDriverRepo(ctrl)
	locationClient := mock.NewMockLocationClient(ctrl)
	eventProducer := mock.NewMockEventProducer(ctrl)
	driverService := svc.NewDriverService(&testDriverServiceConfig, driverRepo, locationClient, eventProducer)
	return driverService, driverRepo, locationClient, eventProducer
}

func TestGetTrips(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrips(gomock.Any(), TEST_DRIVER_ID).Return(testTrips)
	assert.ElementsMatch(t, testTrips, driverService.GetTrips(context.Background(), TEST_DRIVER_ID))
}

func TestGetTrip(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	cases := []struct {
		input    string
		expected *model.Trip
	}{
		{testTrips[0].Id, &testTrips[0]},
		{testTrips[1].Id, &testTrips[1]},
		{TEST_NON_EXISTING_TRIP_ID, nil},
	}
	for _, testCase := range cases {
		driverRepo.EXPECT().GetTrip(gomock.Any(), testCase.input, TEST_DRIVER_ID).Return(testCase.expected)
		assert.Equal(t, testCase.expected, driverService.GetTrip(context.Background(), TEST_DRIVER_ID, testCase.input))
	}
}

func TestCancelTripNotFound(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), TEST_NON_EXISTING_TRIP_ID, TEST_DRIVER_ID).Return(nil)
	assert.Equal(t, svc.TripNotFound, driverService.CancelTrip(context.Background(), TEST_DRIVER_ID, TEST_NON_EXISTING_TRIP_ID, TEST_CANCEL_REASON))
}

func TestCancelTripInvalidChangeStatus(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[4].Id, testTrips[4].DriverId).Return(&testTrips[4])
	assert.Equal(t, svc.InvalidNewStatus, driverService.CancelTrip(context.Background(), testTrips[4].DriverId, testTrips[4].Id, TEST_CANCEL_REASON))
}

func TestCancelTripRepoPutStatusEx(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[0].Id, testTrips[0].DriverId).Return(&testTrips[0])
	driverRepo.EXPECT().PutTripStatus(gomock.Any(), testTrips[0].Id, testTrips[0].DriverId, model.CANCELED).Return(false)
	assert.Equal(t, svc.InternalError, driverService.CancelTrip(context.Background(), testTrips[0].DriverId, testTrips[0].Id, TEST_CANCEL_REASON))
}

func TestCancelTripSuccess(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[0].Id, testTrips[0].DriverId).Return(&testTrips[0])
	driverRepo.EXPECT().PutTripStatus(gomock.Any(), testTrips[0].Id, testTrips[0].DriverId, model.CANCELED).Return(true)
	driverRepo.EXPECT().SaveCancelReason(gomock.Any(), testTrips[0].Id, testTrips[0].DriverId, TEST_CANCEL_REASON).Do(func(_ context.Context, _ string, _ string, _ string) { wg.Done() })
	assert.Equal(t, svc.Ok, driverService.CancelTrip(context.Background(), testTrips[0].DriverId, testTrips[0].Id, TEST_CANCEL_REASON))
	go func() {
		time.Sleep(time.Second * 3)
		t.Error("SaveCancelReason didn't call")
		panic("")
	}()
	wg.Wait()
}

func TestAcceptTripNotFound(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), TEST_NON_EXISTING_TRIP_ID, TEST_DRIVER_ID).Return(nil)
	assert.Equal(t, svc.TripNotFound, driverService.AcceptTrip(context.Background(), TEST_DRIVER_ID, TEST_NON_EXISTING_TRIP_ID))
}

func TestAcceptTripInvalidChangeStatus(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[4].Id, testTrips[4].DriverId).Return(&testTrips[4])
	assert.Equal(t, svc.InvalidNewStatus, driverService.AcceptTrip(context.Background(), testTrips[4].DriverId, testTrips[4].Id))
}

func TestAcceptTripRepoPutStatusEx(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[0].Id, testTrips[0].DriverId).Return(&testTrips[0])
	driverRepo.EXPECT().PutTripStatus(gomock.Any(), testTrips[0].Id, testTrips[0].DriverId, model.DRIVER_FOUND).Return(false)
	assert.Equal(t, svc.InternalError, driverService.AcceptTrip(context.Background(), testTrips[0].DriverId, testTrips[0].Id))
}

func TestAcceptTripSuccess(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	driverService, driverRepo, _, eventProducer := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[0].Id, testTrips[0].DriverId).Return(&testTrips[0])
	driverRepo.EXPECT().PutTripStatus(gomock.Any(), testTrips[0].Id, testTrips[0].DriverId, model.DRIVER_FOUND).Return(true)
	eventProducer.EXPECT().SendTripEvent(gomock.Any(), gomock.Any()).Do(func(_ context.Context, _ *svc.TripMessagePayload) { wg.Done() })
	assert.Equal(t, svc.Ok, driverService.AcceptTrip(context.Background(), testTrips[0].DriverId, testTrips[0].Id))
	go func() {
		time.Sleep(time.Second * 3)
		t.Error("SaveCancelReason didn't call")
		panic("")
	}()
	wg.Wait()
}

func TestStartTripNotFound(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), TEST_NON_EXISTING_TRIP_ID, TEST_DRIVER_ID).Return(nil)
	assert.Equal(t, svc.TripNotFound, driverService.StartTrip(context.Background(), TEST_DRIVER_ID, TEST_NON_EXISTING_TRIP_ID))
}

func TestStartTripInvalidChangeStatus(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[4].Id, testTrips[4].DriverId).Return(&testTrips[4])
	assert.Equal(t, svc.InvalidNewStatus, driverService.StartTrip(context.Background(), testTrips[4].DriverId, testTrips[4].Id))
}

func TestStartTripRepoPutStatusEx(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[2].Id, testTrips[2].DriverId).Return(&testTrips[2])
	driverRepo.EXPECT().PutTripStatus(gomock.Any(), testTrips[2].Id, testTrips[2].DriverId, model.STARTED).Return(false)
	assert.Equal(t, svc.InternalError, driverService.StartTrip(context.Background(), testTrips[2].DriverId, testTrips[2].Id))
}

func TestStartTripSuccess(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	driverService, driverRepo, _, eventProducer := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[2].Id, testTrips[2].DriverId).Return(&testTrips[2])
	driverRepo.EXPECT().PutTripStatus(gomock.Any(), testTrips[2].Id, testTrips[2].DriverId, model.STARTED).Return(true)
	eventProducer.EXPECT().SendTripEvent(gomock.Any(), gomock.Any()).Do(func(_ context.Context, _ *svc.TripMessagePayload) { wg.Done() })
	assert.Equal(t, svc.Ok, driverService.StartTrip(context.Background(), testTrips[2].DriverId, testTrips[2].Id))
	go func() {
		time.Sleep(time.Second * 3)
		t.Error("SaveCancelReason didn't call")
		panic("")
	}()
	wg.Wait()
}

func TestEndTripNotFound(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), TEST_NON_EXISTING_TRIP_ID, TEST_DRIVER_ID).Return(nil)
	assert.Equal(t, svc.TripNotFound, driverService.EndTrip(context.Background(), TEST_DRIVER_ID, TEST_NON_EXISTING_TRIP_ID))
}

func TestEndTripInvalidChangeStatus(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[4].Id, testTrips[4].DriverId).Return(&testTrips[4])
	assert.Equal(t, svc.InvalidNewStatus, driverService.EndTrip(context.Background(), testTrips[4].DriverId, testTrips[4].Id))
}

func TestEndTripRepoPutStatusEx(t *testing.T) {
	driverService, driverRepo, _, _ := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[3].Id, testTrips[3].DriverId).Return(&testTrips[3])
	driverRepo.EXPECT().PutTripStatus(gomock.Any(), testTrips[3].Id, testTrips[3].DriverId, model.ENDED).Return(false)
	assert.Equal(t, svc.InternalError, driverService.EndTrip(context.Background(), testTrips[3].DriverId, testTrips[3].Id))
}

func TestEndTripSuccess(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	driverService, driverRepo, _, eventProducer := initMocksAndDriverService(t)
	driverRepo.EXPECT().GetTrip(gomock.Any(), testTrips[3].Id, testTrips[3].DriverId).Return(&testTrips[3])
	driverRepo.EXPECT().PutTripStatus(gomock.Any(), testTrips[3].Id, testTrips[3].DriverId, model.ENDED).Return(true)
	eventProducer.EXPECT().SendTripEvent(gomock.Any(), gomock.Any()).Do(func(_ context.Context, _ *svc.TripMessagePayload) { wg.Done() })
	assert.Equal(t, svc.Ok, driverService.EndTrip(context.Background(), testTrips[3].DriverId, testTrips[3].Id))
	go func() {
		time.Sleep(time.Second * 3)
		t.Error("SaveCancelReason didn't call")
		panic("")
	}()
	wg.Wait()
}

func TestGetDriversPayloadJsonSerialize(t *testing.T) {
	tmp := svc.GetDriversPayload{
		Radius: 1,
		LatLngLiteral: model.LatLngLiteral{
			Lat: 2,
			Lng: 3,
		},
	}
	a, _ := json.Marshal(tmp)
	assert.Equal(t, string(a), "{\"lat\":2,\"lng\":3,\"radius\":1}")
}

func TestGetDriversPayloadJsonDeserialize(t *testing.T) {
	s := "{\"lat\":2,\"lng\":3,\"radius\":1}"
	var tmp svc.GetDriversPayload
	_ = json.Unmarshal([]byte(s), &tmp)
	assert.Equal(t, 1., tmp.Radius)
	assert.Equal(t, 3., tmp.Lng)
	assert.Equal(t, 2., tmp.Lat)
}
