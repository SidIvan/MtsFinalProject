package httpadapter

import "errors"

type Error struct {
	Message string `json:"message" example:"not"`
}

var (
	ErrNotFound = errors.New("not found")
)
