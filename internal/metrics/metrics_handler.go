package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	NumCallsGetTrips   prometheus.Counter
	TimeProcGetTrips   prometheus.Summary
	NumCallsGetTrip    prometheus.Counter
	TimeProcGetTrip    prometheus.Summary
	NumCallsCancelTrip prometheus.Counter
	TimeProcCancelTrip prometheus.Summary
	NumCallsAcceptTrip prometheus.Counter
	TimeProcAcceptTrip prometheus.Summary
	NumCallsEndTrip    prometheus.Counter
	TimeProcEndTrip    prometheus.Summary
	NumCallsStartTrip  prometheus.Counter
	TimeProcStartTrip  prometheus.Summary
}

func NewMetrics() *Metrics {
	return &Metrics{
		NumCallsGetTrips: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "http_route",
			Name:      "number_of_calls_get_trips",
		}),
		TimeProcGetTrips: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace:  "http_route",
			Name:       "time_of_processing_get_trips_ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		NumCallsGetTrip: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "http_route",
			Name:      "number_of_calls_get_trip",
		}),
		TimeProcGetTrip: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace:  "http_route",
			Name:       "time_of_processing_get_trip_ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		NumCallsCancelTrip: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "http_route",
			Name:      "number_of_calls_cancel_trip",
		}),
		TimeProcCancelTrip: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace:  "http_route",
			Name:       "time_of_processing_cancel_trip_ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		NumCallsAcceptTrip: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "http_route",
			Name:      "number_of_calls_accept_trip",
		}),
		TimeProcAcceptTrip: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace:  "http_route",
			Name:       "time_of_processing_accept_trip_ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		NumCallsEndTrip: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "http_route",
			Name:      "number_of_calls_end_trip",
		}),
		TimeProcEndTrip: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace:  "http_route",
			Name:       "time_of_processing_end_trip_ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		NumCallsStartTrip: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "http_route",
			Name:      "number_of_calls_start_trip",
		}),
		TimeProcStartTrip: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace:  "http_route",
			Name:       "time_of_processing_start_trip_ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
	}
}
