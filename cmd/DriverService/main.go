package main

import (
	"context"
	"driver-service/internal/app"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/otel"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	shutdown := otel.InitProvider()
	defer shutdown()
	cfg := config.NewDriverServiceConfigFromEnv()
	err := logger.InitMainLogger(cfg.EnvType == "dev")
	if err != nil {
		panic(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	a := app.NewDriverServiceApp(cfg)
	a.Start()
	<-ctx.Done()
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(cfg.GracefulShutdownTimeoutSec)*time.Second)
	a.Stop(ctx)
}
