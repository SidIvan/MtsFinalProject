package repomock

import (
	"context"

	"gitlab.com/AntYats/go_project/internal/model"
	"gitlab.com/AntYats/go_project/internal/repo"

	"github.com/stretchr/testify/mock"
)

var _ repo.User = &UserMock{}

type UserMock struct {
	mock.Mock
}

func (m *UserMock) WithNewTx(ctx context.Context, f func(ctx context.Context) error) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *UserMock) ChangeDriverInfo(ctx context.Context, id string, coords *model.UserData) error {
	args := m.Called(ctx, id, coords)
	return args.Error(0)
}

func (m *UserMock) GetDrivers(ctx context.Context, radius float64, coords *model.UserData) ([]model.User, error) {
	args := m.Called(ctx, radius, coords)
	return args.Get(0).([]model.User), args.Error(1)
}

func NewUser() *UserMock {
	return &UserMock{}
}
