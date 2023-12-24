package repo

import (
	"context"

	"gitlab.com/AntYats/go_project/internal/model"
)

type User interface {
	GetDrivers(ctx context.Context, radius float64, coords *model.UserData) ([]model.User, error)
	ChangeDriverInfo(ctx context.Context, id string, coords *model.UserData) error
	WithNewTx(ctx context.Context, f func(ctx context.Context) error) error
}
