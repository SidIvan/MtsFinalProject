package repo

import (
	"context"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/model"
	"driver-service/internal/repo"
	"driver-service/internal/svc"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

const (
	TEST_TIMEOUT    = 3
	NON_EXISTING_ID = "non_existing_id"
	TEST_REASON     = "test_reason"
)

var (
	cont              *mongodb.MongoDBContainer
	driverRepo        *repo.DriverRepoMongo
	tripColl          *mongo.Collection
	cancelLogColl     *mongo.Collection
	TEST_BASE_CONTEXT = context.Background()
	testTrip          = &model.Trip{
		Id:       "test_id",
		DriverId: "driver_id",
		From: model.LatLngLiteral{
			Lat: 1.,
			Lng: 2.,
		},
		To: model.LatLngLiteral{
			Lat: 3.,
			Lng: 4.,
		},
		Price: model.Money{
			Currency: "MOP",
			Amount:   5,
		},
		Status: model.STARTED,
	}
	testTrip2 = &model.Trip{
		Id:       "test_id_2",
		DriverId: "driver_id_2",
		From: model.LatLngLiteral{
			Lat: 10.,
			Lng: 20.,
		},
		To: model.LatLngLiteral{
			Lat: 30.,
			Lng: 40.,
		},
		Price: model.Money{
			Currency: "RUB",
			Amount:   50,
		},
		Status: model.ENDED,
	}
	testTrip3 = &model.Trip{
		Id:       "test_id_3",
		DriverId: "driver_id",
		From: model.LatLngLiteral{
			Lat: 100.,
			Lng: 200.,
		},
		To: model.LatLngLiteral{
			Lat: 300.,
			Lng: 400.,
		},
		Price: model.Money{
			Currency: "USD",
			Amount:   500,
		},
		Status: model.ON_POSITION,
	}
)

func initMongoTestContainer(t *testing.T) {
	var err error
	cont, err = mongodb.RunContainer(context.Background(), testcontainers.WithImage("mongo:6"))
	if err != nil {
		t.Error(err.Error())
	}
}

func beforeTestPipeline(t *testing.T) {
	logger.InitTestLogger(t)
	initMongoTestContainer(t)
	connStr, err := cont.ConnectionString(context.Background())
	if err != nil {
		t.Error(err.Error())
	}
	driverRepo = repo.NewDriverRepoMongo(&config.MongoConfig{
		TimeoutSec:                    TEST_TIMEOUT,
		URI:                           connStr,
		DatabaseName:                  "test",
		TripCollectionName:            "test_trips",
		CancelReasonLogCollectionName: "test_cancel_reasons",
	})
	tripColl = driverRepo.GetTripCollection()
	cancelLogColl = driverRepo.GetCancelLogCollection()
}

func afterTestPipeline(t *testing.T) {
	err := cont.Terminate(context.Background())
	if err != nil {
		t.Error(err.Error())
	}
}

func TestCreateTripSuccess(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	defer tripColl.Database().Client().Disconnect(context.Background())
	id := driverRepo.CreateTrip(TEST_BASE_CONTEXT, testTrip)
	var savedTrip model.Trip
	res := tripColl.FindOne(context.Background(), bson.M{"_id": id})
	err := res.Decode(&savedTrip)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, testTrip.Id, id)
	assert.Equal(t, *testTrip, savedTrip)
}

func TestCreateTripFailed(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	defer tripColl.Database().Client().Disconnect(context.Background())
	driverRepo.CreateTrip(TEST_BASE_CONTEXT, testTrip)
	id := driverRepo.CreateTrip(TEST_BASE_CONTEXT, testTrip)
	assert.Equal(t, "", id)
}

func TestPutTripStatusSuccess(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	defer tripColl.Database().Client().Disconnect(context.Background())
	putStatus := model.CANCELED
	expectedRes := *testTrip
	expectedRes.Status = putStatus
	id := driverRepo.CreateTrip(TEST_BASE_CONTEXT, testTrip)
	if !driverRepo.PutTripStatus(context.Background(), id, testTrip.DriverId, putStatus) {
		t.Errorf("Status was not change")
		return
	}
	var savedTrip model.Trip
	res := tripColl.FindOne(context.Background(), bson.M{"_id": id})
	err := res.Decode(&savedTrip)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, expectedRes, savedTrip)
}

func TestPutTripStatusInvalidTripIdFailed(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	defer tripColl.Database().Client().Disconnect(context.Background())
	putStatus := model.CANCELED
	expected_res := *testTrip
	expected_res.Status = putStatus
	driverRepo.CreateTrip(TEST_BASE_CONTEXT, testTrip)
	assert.Equal(t, false, driverRepo.PutTripStatus(context.Background(), NON_EXISTING_ID, testTrip.DriverId, putStatus))
}

func TestPutTripStatusConnErrorFailed(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	putStatus := model.CANCELED
	expected_res := *testTrip
	expected_res.Status = putStatus
	driverRepo.CreateTrip(TEST_BASE_CONTEXT, testTrip)
	tripColl.Database().Client().Disconnect(context.Background())
	assert.Equal(t, false, driverRepo.PutTripStatus(context.Background(), NON_EXISTING_ID, testTrip.DriverId, putStatus))
}

func TestGetTripSuccessfull(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	_, err := tripColl.InsertOne(TEST_BASE_CONTEXT, testTrip)
	if err != nil {
		t.Error(err.Error())
		return
	}
	savedTrip := driverRepo.GetTrip(TEST_BASE_CONTEXT, testTrip.Id, testTrip.DriverId)
	assert.Equal(t, testTrip, savedTrip)
}

func TestGetTripInvalidIdFalied(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	_, err := tripColl.InsertOne(TEST_BASE_CONTEXT, testTrip)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Nil(t, driverRepo.GetTrip(TEST_BASE_CONTEXT, NON_EXISTING_ID, testTrip.DriverId))
}

func TestGetTripConnExFalied(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	_, err := tripColl.InsertOne(TEST_BASE_CONTEXT, testTrip)
	tripColl.Database().Client().Disconnect(context.Background())
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Nil(t, driverRepo.GetTrip(TEST_BASE_CONTEXT, NON_EXISTING_ID, testTrip.DriverId))
}

func fillTripCollection(t *testing.T) {
	for _, trip := range []*model.Trip{testTrip, testTrip2, testTrip3} {
		_, err := tripColl.InsertOne(TEST_BASE_CONTEXT, trip)
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
}

func TestGetTripsSuccess(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	fillTripCollection(t)
	result := driverRepo.GetTrips(TEST_BASE_CONTEXT, testTrip.DriverId)
	assert.ElementsMatch(t, []model.Trip{*testTrip, *testTrip3}, result)
}

func TestGetTripsConnExFailed(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	fillTripCollection(t)
	tripColl.Database().Client().Disconnect(context.Background())
	result := driverRepo.GetTrips(TEST_BASE_CONTEXT, testTrip.DriverId)
	assert.ElementsMatch(t, nil, result)
}

func TestSaveCancelReasonSuccess(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	expectedMessage := svc.NewCancelReasonLog(testTrip.Id, testTrip.DriverId, TEST_REASON)
	driverRepo.SaveCancelReason(TEST_BASE_CONTEXT, testTrip.Id, testTrip.DriverId, TEST_REASON)
	var actualMessage svc.CancelReasonLog
	result := cancelLogColl.FindOne(TEST_BASE_CONTEXT, bson.M{"_id": testTrip.Id})
	err := result.Decode(&actualMessage)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, expectedMessage, &actualMessage)
}

func TestSaveCancelReasonConnExFailed(t *testing.T) {
	beforeTestPipeline(t)
	defer afterTestPipeline(t)
	tripColl.Database().Client().Disconnect(context.Background())
	driverRepo.SaveCancelReason(TEST_BASE_CONTEXT, testTrip.Id, testTrip.DriverId, TEST_REASON)
}
