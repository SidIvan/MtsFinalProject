package model

type Trip struct {
	Id       string        `bson:"_id" json:"id"`
	DriverId string        `bson:"driver_id" json:"driver_id"`
	From     LatLngLiteral `bson:"from" json:"from"`
	To       LatLngLiteral `bson:"to" json:"to"`
	Price    Money         `bson:"money" json:"money"`
	Status   TripStatus    `bson:"status" json:"status"`
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

func (t *Trip) IsValid() bool {
	return t.From.isValid() && t.To.isValid() && t.Price.IsValid() &&
		(t.Status == CREATED || t.Status == DRIVER_FOUND || t.Status == ON_POSITION ||
			t.Status == STARTED || t.Status == ENDED || t.Status == CANCELED)
}

func IsValidChangeStatus(from TripStatus, to TripStatus) bool {
	return (from == CREATED && to == DRIVER_FOUND) ||
		(from == DRIVER_FOUND && to == ON_POSITION) ||
		(from == ON_POSITION && to == STARTED) ||
		(from == STARTED && to == ENDED) ||
		(from != ENDED && to == CANCELED)
}
