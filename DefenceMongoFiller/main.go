package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
)

type LatLngLiteral struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lng float64 `bson:"lng" json:"lng"`
}

type TripStatus string

const (
	CREATED      = TripStatus("CREATED")
	DRIVER_FOUND = TripStatus("DRIVER_FOUND")
	ON_POSITION  = TripStatus("ON_POSITION")
	STARTED      = TripStatus("STARTED")
	ENDED        = TripStatus("ENDED")
	CANCELED     = TripStatus("CANCELED")
)

type Trip struct {
	Id       string        `bson:"_id" json:"id"`
	DriverId string        `bson:"driver_id" json:"driver_id"`
	From     LatLngLiteral `bson:"from" json:"from"`
	To       LatLngLiteral `bson:"to" json:"to"`
	Price    Money         `bson:"money" json:"money"`
	Status   TripStatus    `bson:"status" json:"status"`
}

var statuses = []TripStatus{
	CREATED, DRIVER_FOUND, ON_POSITION, STARTED, ENDED, CANCELED,
}

type Money struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
}

func main() {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	tripCol := client.Database("database").Collection("trips")
	//clCol := client.Database("database").Collection("cancel_logs")
	for i := 0; i < 10; i++ {
		tripCol.InsertOne(context.TODO(), Trip{
			Id:       "id" + strconv.Itoa(i),
			DriverId: "driverId" + strconv.Itoa(i),
			From: LatLngLiteral{
				Lat: float64(i),
				Lng: float64(i),
			},
			To: LatLngLiteral{
				Lat: float64(i) + 0.5,
				Lng: float64(i) + 0.5,
			},
			Price: Money{
				Amount:   float64(i),
				Currency: "MOP",
			},
			Status: statuses[i%6],
		})
	}
	for i := 0; i < 10; i++ {
		tripCol.InsertOne(context.TODO(), Trip{
			Id:       "id1" + strconv.Itoa(i),
			DriverId: "driverId" + strconv.Itoa(i),
			From: LatLngLiteral{
				Lat: float64(i),
				Lng: float64(i),
			},
			To: LatLngLiteral{
				Lat: float64(i) + 0.5,
				Lng: float64(i) + 0.5,
			},
			Price: Money{
				Amount:   float64(i),
				Currency: "RUB",
			},
			Status: statuses[(i*5)%6],
		})
	}
	var result Trip
	filter := bson.M{"_id": "id18", "driver_id": "driverId8"}
	err := tripCol.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Found document:", result)

}
