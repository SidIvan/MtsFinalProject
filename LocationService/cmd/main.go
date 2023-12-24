package main

import (
	"flag"
	"gitlab.com/AntYats/go_project/internal/app"
	"go.uber.org/zap"
)

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "../../.config/location.env.dev", "path to config file")
	flag.Parse()

	return configPath
}

func main() {
	logger, err := zap.NewProduction()

	if err != nil {
		logger.Error("Logger initialization error", zap.Error(err))
	}

	config, err := app.NewConfig(getConfigPath())

	if err != nil {
		logger.Error("Load configuration error", zap.Error(err))
	}

	a, err := app.New(config, logger)
	if err != nil {
		logger.Error("App creation error", zap.Error(err))
	}

	if err := a.Serve(); err != nil {
		logger.Error("App error", zap.Error(err))
	}
}
