package web

import (
	"bytes"
	"context"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/model"
	"driver-service/internal/svc"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LocationClientHttp struct {
	client   *http.Client
	basePath string
	timeout  time.Duration
}

func NewLocationClientHttp(cfg *config.LocationServiceConfig) *LocationClientHttp {
	return &LocationClientHttp{
		client:   &http.Client{},
		timeout:  time.Duration(cfg.TimeoutSec) * time.Second,
		basePath: fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port),
	}
}

func (c LocationClientHttp) GetDrivers(baseCtx context.Context, payload *svc.GetDriversPayload) []model.Driver {
	ctx, cancel := context.WithTimeout(baseCtx, c.timeout)
	defer cancel()
	body, err := json.Marshal(payload)
	if err != nil {
		logger.Main.Error(err.Error())
		return nil
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/drivers", c.basePath), bytes.NewReader(body))
	logger.Main.Info(fmt.Sprintf("Sending http request %s %s", request.Method, request.URL.Path))
	response, err := c.client.Do(request)
	if err != nil {
		logger.Main.Error(fmt.Sprintf("Failed http request %s %s", request.Method, request.URL.Path))
		return nil
	}
	if response.Status == "404 Not found" {
		logger.Main.Warn(fmt.Sprintf("Empty response %s %s", request.Method, request.URL.Path))
		return nil
	}
	logger.Main.Info(fmt.Sprintf("Success http request %s %s", request.Method, request.URL.Path))
	responseBody, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		logger.Main.Error(err.Error())
		return nil
	}
	var drivers []model.Driver
	err = json.Unmarshal(responseBody, &drivers)
	if err != nil {
		logger.Main.Error(err.Error())
		return nil
	}
	return drivers
}
