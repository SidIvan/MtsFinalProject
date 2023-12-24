package service

import (
	"context"
	"gitlab.com/AntYats/go_project/internal/model"
)

type User interface {
	GetDrivers(ctx context.Context, radius float64, coords *model.UserData) ([]model.User, error)
	ChangeDriverInfo(ctx context.Context, id string, coords *model.UserData) error
}
