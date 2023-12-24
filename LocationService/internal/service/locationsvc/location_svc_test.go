package locationsvc_test

import (
	"context"
	"testing"

	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gitlab.com/AntYats/go_project/internal/model"
	// "gitlab.com/AntYats/go_project/internal/repo"
	"gitlab.com/AntYats/go_project/internal/repo/repomock"
	// "gitlab.com/AntYats/go_project/internal/service"
	"gitlab.com/AntYats/go_project/internal/service/locationsvc"
)

func TestGetDrivers(t *testing.T) {
	ctx := context.Background()
	users := []model.User{}

	repo := repomock.NewUser()

	serviceMy := locationsvc.New(repo)
	users_res, err := serviceMy.GetDrivers(ctx, 1, nil)
	require.NoError(t, err)
	require.Equal(t, users, users_res)
}

func TestChangeDriverInfo(t *testing.T) {
	ctx := context.Background()

	repo := repomock.NewUser()

	serviceMy := locationsvc.New(repo)
	err := serviceMy.ChangeDriverInfo(ctx, "111", &model.UserData{Lat: 1, Lng: 1})
	require.NoError(t, err)
}
