package model

type User struct {
	UserId string  `json:"user_id" db:"user_id"`
	Lat    float32 `json:"lat"  db:"lat"`
	Lng    float32 `json:"lng"  db:"lng"`
}

type UserData struct {
	Lat float32 `json:"lat"  db:"lat"`
	Lng float32 `json:"lng"  db:"lng"`
}
