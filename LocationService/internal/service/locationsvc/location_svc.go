package locationsvc

import (
	"context"
	"gitlab.com/AntYats/go_project/internal/model"
	"gitlab.com/AntYats/go_project/internal/repo"
)

type locationService struct {
	repo repo.User
}

func (l *locationService) GetDrivers(ctx context.Context, radius float64, coords *model.UserData) ([]model.User, error) {
	return l.repo.GetDrivers(ctx, radius, coords)
}

func (l *locationService) ChangeDriverInfo(ctx context.Context, id string, coords *model.UserData) error {
	return l.repo.WithNewTx(ctx, func(ctx context.Context) error {
		return l.repo.ChangeDriverInfo(ctx, id, coords)
	})
}

func New(repo repo.User) *locationService {
	return &locationService{
		repo: repo,
	}
}
