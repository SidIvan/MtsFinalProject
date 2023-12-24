package svc

type CancelReasonLog struct {
	TripId   string `bson:"_id" json:"trip_id"`
	DriverId string `bson:"driver_id" json:"driver_id"`
	Reason   string `bson:"reason" json:"reason"`
}

func NewCancelReasonLog(tripId string, driverId string, reason string) *CancelReasonLog {
	return &CancelReasonLog{
		TripId:   tripId,
		DriverId: driverId,
		Reason:   reason,
	}
}
