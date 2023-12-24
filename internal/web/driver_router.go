package web

import (
	"driver-service/internal/logger"
	"driver-service/internal/svc"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var (
	USER_ID_HEADER_POSSIBLE_KEYS = map[string]struct{}{
		"user_id": {},
		"User_id": {},
	}
)

type DriverRouter struct {
	driverService svc.DriverService
}

func NewDriverRouter(driverService svc.DriverService) *DriverRouter {
	return &DriverRouter{
		driverService: driverService,
	}
}

func writeResponseBody(w http.ResponseWriter, payload interface{}) {
	if payload == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	responseData, err := json.Marshal(payload)
	if err != nil {
		logger.Main.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writtenBytes, err := w.Write(responseData)
	if err != nil {
		logger.Main.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if writtenBytes != len(responseData) {
		logger.Main.Error("Not full response written")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getUserIdHeaderAnd404IfNotFound(w http.ResponseWriter, r *http.Request) (string, bool) {
	for headerKey, _ := range USER_ID_HEADER_POSSIBLE_KEYS {
		headerVal := r.Header.Get(headerKey)
		if headerVal != "" {
			return headerVal, true
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	return "", false
}

func (rtr *DriverRouter) GetTripsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := svc.Tracer.Start(r.Context(), "/get-trips")
	defer span.End()
	if userId, ok := getUserIdHeaderAnd404IfNotFound(w, r); ok {
		trips := rtr.driverService.GetTrips(ctx, userId)
		writeResponseBody(w, trips)
	}
}

func (rtr *DriverRouter) GetTripHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := svc.Tracer.Start(r.Context(), "/get-trip")
	defer span.End()
	if userId, ok := getUserIdHeaderAnd404IfNotFound(w, r); ok {
		tripId := mux.Vars(r)["trip_id"]
		trip := rtr.driverService.GetTrip(ctx, userId, tripId)
		if trip == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		writeResponseBody(w, trip)
	}
}

var tripOpStatusToHttpStatus = map[svc.ChangeTripStatusStatus]int{
	svc.Ok:               http.StatusOK,
	svc.InvalidNewStatus: http.StatusBadRequest,
	svc.TripNotFound:     http.StatusNotFound,
	svc.InternalError:    http.StatusInternalServerError,
}

func (rtr *DriverRouter) CancelTripHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := svc.Tracer.Start(r.Context(), "/cancel-trip")
	defer span.End()
	if userId, ok := getUserIdHeaderAnd404IfNotFound(w, r); ok {
		tripId := mux.Vars(r)["trip_id"]
		reason, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logger.Main.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		opRes := rtr.driverService.CancelTrip(ctx, userId, tripId, string(reason))
		w.WriteHeader(tripOpStatusToHttpStatus[opRes])
	}
}

func (rtr *DriverRouter) AcceptTripHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := svc.Tracer.Start(r.Context(), "/accept-trip")
	defer span.End()
	if userId, ok := getUserIdHeaderAnd404IfNotFound(w, r); ok {
		tripId := mux.Vars(r)["trip_id"]
		opRes := rtr.driverService.AcceptTrip(ctx, userId, tripId)
		w.WriteHeader(tripOpStatusToHttpStatus[opRes])
	}
}

func (rtr *DriverRouter) StartTripHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := svc.Tracer.Start(r.Context(), "/start-trip")
	defer span.End()
	if userId, ok := getUserIdHeaderAnd404IfNotFound(w, r); ok {
		tripId := mux.Vars(r)["trip_id"]
		opRes := rtr.driverService.StartTrip(ctx, userId, tripId)
		w.WriteHeader(tripOpStatusToHttpStatus[opRes])
	}
}

func (rtr *DriverRouter) EndTripHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := svc.Tracer.Start(r.Context(), "/end-trip")
	defer span.End()
	if userId, ok := getUserIdHeaderAnd404IfNotFound(w, r); ok {
		tripId := mux.Vars(r)["trip_id"]
		opRes := rtr.driverService.EndTrip(ctx, userId, tripId)
		w.WriteHeader(tripOpStatusToHttpStatus[opRes])
	}
}
