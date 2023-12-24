package svc

import "driver-service/internal/model"

type GetDriversPayload struct {
	model.LatLngLiteral
	Radius float64 `json:"radius"`
}
