package app

import (
	"context"
	"driver-service/internal/config"
	"driver-service/internal/logger"
	"driver-service/internal/metrics"
	"driver-service/internal/repo"
	"driver-service/internal/svc"
	"driver-service/internal/web"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type App struct {
	httpServer    *http.Server
	eventConsumer *web.CreateEventConsumer
}

func NewDriverServiceApp(cfg *config.DriverServiceConfig) *App {
	locationClient := web.NewLocationClientHttp(cfg.LocationServiceConfig)
	eventProducer := web.NewKafkaEventProducer(cfg.KafkaConfig)
	driverRepo := repo.NewDriverRepoMongo(cfg.MongoConfig)
	driverService := svc.NewDriverService(cfg, driverRepo, locationClient, eventProducer)
	initNat(cfg.NatConfig)
	apiMetrics := metrics.NewMetrics()
	return &App{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: initApi(driverService, apiMetrics),
		},
		eventConsumer: web.NewCreateEventConsumer(cfg.KafkaConfig, *driverService),
	}
}

func (a *App) Start() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":9001", nil)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		for {
			a.eventConsumer.Start(context.Background())
		}
	}()
	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	logger.Main.Info(fmt.Sprintf("Driver service started at port %s", a.httpServer.Addr))
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "example_counter",
		Help: "This is an example counter",
	})
	counter.Inc()
	counter.Inc()
	counter.Inc()
}

func (a *App) Stop(ctx context.Context) {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		logger.Main.Info("Driver service shutdown gracefully")
	}
	logger.Main.Error("Error while gracefull shutdown")
	a.eventConsumer.Stop()
}

func initApi(driverService *svc.DriverServiceImpl, apiMetrics *metrics.Metrics) http.Handler {
	controller := web.NewDriverRouter(driverService)
	r := mux.NewRouter()
	r.
		HandleFunc("/trips", func(w http.ResponseWriter, r *http.Request) {
			apiMetrics.NumCallsGetTrips.Inc()
			startTime := time.Now()
			controller.GetTripsHandler(w, r)
			apiMetrics.TimeProcGetTrips.Observe(float64(time.Since(startTime).Milliseconds()))
		}).
		Methods(http.MethodGet)
	r.
		HandleFunc("/trip/{trip_id}", func(w http.ResponseWriter, r *http.Request) {
			apiMetrics.NumCallsGetTrip.Inc()
			startTime := time.Now()
			controller.GetTripHandler(w, r)
			apiMetrics.TimeProcGetTrip.Observe(float64(time.Since(startTime).Milliseconds()))
		}).
		Methods(http.MethodGet)
	r.
		HandleFunc("/trip/{trip_id}/cancel", func(w http.ResponseWriter, r *http.Request) {
			apiMetrics.NumCallsCancelTrip.Inc()
			startTime := time.Now()
			controller.CancelTripHandler(w, r)
			apiMetrics.TimeProcCancelTrip.Observe(float64(time.Since(startTime).Milliseconds()))
		}).
		Methods(http.MethodPost)
	r.
		HandleFunc("/trip/{trip_id}/accept", func(w http.ResponseWriter, r *http.Request) {
			apiMetrics.NumCallsAcceptTrip.Inc()
			startTime := time.Now()
			controller.AcceptTripHandler(w, r)
			apiMetrics.TimeProcAcceptTrip.Observe(float64(time.Since(startTime).Milliseconds()))
		}).
		Methods(http.MethodPost)
	r.
		HandleFunc("/trip/{trip_id}/start", func(w http.ResponseWriter, r *http.Request) {
			apiMetrics.NumCallsStartTrip.Inc()
			startTime := time.Now()
			controller.StartTripHandler(w, r)
			apiMetrics.TimeProcStartTrip.Observe(float64(time.Since(startTime).Milliseconds()))
		}).
		Methods(http.MethodPost)
	r.
		HandleFunc("/trip/{trip_id}/end", func(w http.ResponseWriter, r *http.Request) {
			apiMetrics.NumCallsEndTrip.Inc()
			startTime := time.Now()
			controller.EndTripHandler(w, r)
			apiMetrics.TimeProcEndTrip.Observe(float64(time.Since(startTime).Milliseconds()))
		}).
		Methods(http.MethodPost)
	r.Use(logger.MiddleWare)
	return r
}

func initNat(cfg *config.NatNotificationerConfig) {
}
