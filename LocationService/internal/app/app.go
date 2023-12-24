package app

import (
	"context"
	_ "embed"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	// "database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	// "gitlab.com/AntYats/go_project/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"gitlab.com/AntYats/go_project/internal/httpadapter"
	"gitlab.com/AntYats/go_project/internal/repo"
	"gitlab.com/AntYats/go_project/internal/repo/user_repo"
	"gitlab.com/AntYats/go_project/internal/service"
	"gitlab.com/AntYats/go_project/internal/service/locationsvc"
)

type App struct {
	config           *Config
	lotcationService service.User
	httpAdapter      httpadapter.Adapter
	store            repo.User
	logger           *zap.Logger
}

func New(config *Config, logger *zap.Logger) (*App, error) {
	logger.Info("New app process started")
	pgxPool, err := initDB(context.Background(), &config.Database, logger)

	if err != nil {
		logger.Error("Database initialization error")
		return nil, err
	}

	userRepo, err := userrepo.New(pgxPool)
	if err != nil {
		logger.Error("User repo initialization error")
		return nil, err
	}

	locationService := locationsvc.New(userRepo)

	a := &App{
		config:           config,
		lotcationService: locationService,
		httpAdapter:      httpadapter.New(&config.HTTP, locationService),
		store:            userRepo,
	}

	logger.Info("New app process succesfully ended")

	return a, nil
}

func initDB(ctx context.Context, config *DatabaseConfig, logger *zap.Logger) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		logger.Error("Parsing error in db", zap.Error(err))
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		logger.Error("Unable to connect to database", zap.Error(err))
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	m, err := migrate.New(config.MigrationsDir, config.DSN)

	if err != nil {
		logger.Error("Migration error", zap.Error(err))
		return nil, err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		logger.Error("Migration: down - error", zap.Error(err))
		return nil, err
	}

	if err := m.Up(); err != nil {
		logger.Error("Migration: up - error", zap.Error(err))
		return nil, err
	}

	return pool, nil
}

func (a *App) Serve() error {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := a.httpAdapter.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	a.logger.Info("App started.")

	<-done

	a.Shutdown()

	return nil
}

func (a *App) Shutdown() {
	a.logger.Info("App shutting down.")

	ctx, cancel := context.WithTimeout(context.Background(), a.config.App.ShutdownTimeout)
	defer cancel()

	a.httpAdapter.Shutdown(ctx)

	a.logger.Info("App shutdown.")
}
