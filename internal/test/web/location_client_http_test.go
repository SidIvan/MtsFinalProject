package web

import (
	"context"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/model"
	"driver-service/internal/svc"
	"driver-service/internal/web"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"strconv"
	"testing"
	"time"
)

const (
	IMAGE_NAME_200 = "locsvcmock:200"
	IMAGE_NAME_404 = "locsvcmock:404"
	TEST_PORT      = "8081/tcp"
)

var (
	locationServiceTestConfig = &config.LocationServiceConfig{
		Host:       "localhost",
		Port:       0,
		TimeoutSec: 3,
	}
	locationClient *web.LocationClientHttp
)

func initLocationClientMocks(t *testing.T, imageName string) func(context.Context) error {
	logger.InitTestLogger(t)
	contReq := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: imageName,
		},
		Started: true,
	}
	cont, err := testcontainers.GenericContainer(context.Background(), contReq)
	if err != nil {
		t.Error(err.Error())
	}
	var port int
	for i := 0; i < 3; i++ {
		ports, err := cont.Ports(context.TODO())
		if err == nil && len(ports[TEST_PORT]) > 0 {
			port, err = strconv.Atoi(ports[TEST_PORT][0].HostPort)
			if err != nil {
				t.Error(err.Error())
			}
			locationServiceTestConfig.Port = port
			locationClient = web.NewLocationClientHttp(locationServiceTestConfig)
			return cont.Terminate
		}
		time.Sleep(time.Second)
	}
	t.Error("Error init mocks")
	return nil
}

var testPayload = svc.GetDriversPayload{
	LatLngLiteral: model.LatLngLiteral{
		Lat: 42.,
		Lng: 13.,
	},
	Radius: 2,
}

func TestGetDriversSuccess(t *testing.T) {
	defer initLocationClientMocks(t, IMAGE_NAME_200)(context.Background())
	resp := locationClient.GetDrivers(context.Background(), &testPayload)
	assert.True(t, len(resp) > 0)
}

func TestGetDriversNotFoundFailed(t *testing.T) {
	defer initLocationClientMocks(t, IMAGE_NAME_404)(context.Background())
	resp := locationClient.GetDrivers(context.Background(), &testPayload)
	assert.Nil(t, resp)
}

func TestGetDriversConnDownFailed(t *testing.T) {
	logger.InitTestLogger(t)
	locationClient := web.NewLocationClientHttp(locationServiceTestConfig)
	resp := locationClient.GetDrivers(context.Background(), &testPayload)
	assert.Nil(t, resp)
}
