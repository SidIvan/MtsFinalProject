package model

type LatLngLiteral struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lng float64 `bson:"lng" json:"lng"`
}

func (l *LatLngLiteral) isValid() bool {
	return l.IsLongValid() && l.IsLatValid()
}

func (l *LatLngLiteral) IsLatValid() bool {
	return IsLatValid(l.Lat)
}

func (l *LatLngLiteral) IsLongValid() bool {
	return IsLongValid(l.Lat)
}

func IsLatValid(lat float64) bool {
	return lat >= -90 && lat <= 90
}

func IsLongValid(long float64) bool {
	return long > -180 && long <= 180
}
