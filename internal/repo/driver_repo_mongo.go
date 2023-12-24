package repo

import (
	"context"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/model"
	"driver-service/internal/svc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DriverRepoMongo struct {
	tripCollection      *mongo.Collection
	cancelLogCollection *mongo.Collection
	timeout             time.Duration
}

func NewDriverRepoMongo(cfg *config.MongoConfig) *DriverRepoMongo {
	timeout := time.Duration(cfg.TimeoutSec) * time.Second
	database := connectMongo(cfg, timeout)
	return &DriverRepoMongo{
		tripCollection:      database.Collection(cfg.TripCollectionName),
		cancelLogCollection: database.Collection(cfg.CancelReasonLogCollectionName),
		timeout:             timeout,
	}
}

func connectMongo(cfg *config.MongoConfig, timeout time.Duration) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		panic(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	return client.Database(cfg.DatabaseName)
}

func (r DriverRepoMongo) CreateTrip(baseCtx context.Context, trip *model.Trip) string {
	ctx, cancel := context.WithTimeout(baseCtx, r.timeout)
	defer cancel()
	res, err := r.tripCollection.InsertOne(ctx, trip)
	if err != nil {
		logger.Main.Error(err.Error())
		return ""
	}
	return res.InsertedID.(string)
}

func (r DriverRepoMongo) PutTripStatus(baseCtx context.Context, tripId string, driverId string, newStatus model.TripStatus) bool {
	ctx, cancel := context.WithTimeout(baseCtx, r.timeout)
	defer cancel()
	res, err := r.tripCollection.UpdateOne(ctx, bson.M{"_id": tripId, "driver_id": driverId}, bson.M{"$set": bson.M{"status": newStatus}})
	if err != nil {
		logger.Main.Error(err.Error())
		return false
	}
	if res.ModifiedCount != 1 {
		return false
	}
	return true
}

func (r DriverRepoMongo) GetTrip(baseCtx context.Context, tripId string, driverId string) *model.Trip {
	ctx, cancel := context.WithTimeout(baseCtx, r.timeout)
	defer cancel()
	res := r.tripCollection.FindOne(ctx, bson.M{"_id": tripId, "driver_id": driverId})
	if res.Err() != nil {
		logger.Main.Error(res.Err().Error())
		return nil
	}
	var trip model.Trip
	err := res.Decode(&trip)
	if err != nil {
		logger.Main.Error(err.Error())
		return nil
	}
	return &trip
}

func (r DriverRepoMongo) GetTrips(baseCtx context.Context, driverId string) []model.Trip {
	ctx, cancel := context.WithTimeout(baseCtx, r.timeout)
	defer cancel()
	cur, err := r.tripCollection.Find(ctx, bson.M{"driver_id": driverId})
	if err != nil {
		logger.Main.Error(err.Error())
		return nil
	}
	ctx, cancel = context.WithTimeout(baseCtx, r.timeout)
	defer cancel()
	var trips []model.Trip
	err = cur.All(ctx, &trips)
	if err != nil {
		logger.Main.Error(err.Error())
		return nil
	}
	return trips
}

func (r DriverRepoMongo) SaveCancelReason(baseCtx context.Context, tripId string, driverId string, reason string) {
	ctx, cancel := context.WithTimeout(baseCtx, r.timeout)
	defer cancel()
	_, err := r.cancelLogCollection.InsertOne(ctx, svc.NewCancelReasonLog(tripId, driverId, reason))
	if err != nil {
		logger.Main.Error(err.Error())
	}
}

func (r DriverRepoMongo) GetTripCollection() *mongo.Collection {
	return r.tripCollection
}

func (r DriverRepoMongo) GetCancelLogCollection() *mongo.Collection {
	return r.cancelLogCollection
}
