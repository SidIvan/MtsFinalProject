package web

import (
	"bytes"
	"context"
	"driver-service/internal/logger"
	"driver-service/internal/model"
	"driver-service/internal/svc"
	"driver-service/internal/test/mock"
	"driver-service/internal/web"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"testing"
)

const (
	TEST_USER_ID            = "test_user_id"
	TEST_BASE_URI           = "test_uri"
	TEST_TRIP_ID            = "test_trip_id"
	TEST_USER_ID_HEADER_KEY = "user_id"
	TEST_CANCEL_REASON      = "test_cancel_reason"
)

var (
	TEST_BASE_CONTEXT = context.Background()
	driverService     *mock.MockDriverService
	driverRouter      *web.DriverRouter
	respWr            *mock.MockResponseWriter
	trips             = []model.Trip{
		{
			Id:       "test_id",
			DriverId: "driver_id",
			From: model.LatLngLiteral{
				Lat: 1.,
				Lng: 2.,
			},
			To: model.LatLngLiteral{
				Lat: 3.,
				Lng: 4.,
			},
			Price: model.Money{
				Currency: "MOP",
				Amount:   5,
			},
			Status: model.STARTED,
		},
		{
			Id:       "test_id_2",
			DriverId: "driver_id_2",
			From: model.LatLngLiteral{
				Lat: 10.,
				Lng: 20.,
			},
			To: model.LatLngLiteral{
				Lat: 30.,
				Lng: 40.,
			},
			Price: model.Money{
				Currency: "RUB",
				Amount:   50,
			},
			Status: model.ENDED,
		},
		{
			Id:       "test_id_3",
			DriverId: "driver_id",
			From: model.LatLngLiteral{
				Lat: 100.,
				Lng: 200.,
			},
			To: model.LatLngLiteral{
				Lat: 300.,
				Lng: 400.,
			},
			Price: model.Money{
				Currency: "USD",
				Amount:   500,
			},
			Status: model.ON_POSITION,
		},
	}
)

func initDriverRouterMocks(t *testing.T) {
	logger.InitTestLogger(t)
	ctrl := gomock.NewController(t)
	driverService = mock.NewMockDriverService(ctrl)
	driverRouter = web.NewDriverRouter(driverService)
	respWr = mock.NewMockResponseWriter(ctrl)
}

func TestGetTripsHandlerSuccess(t *testing.T) {
	initDriverRouterMocks(t)
	req, err := http.NewRequestWithContext(TEST_BASE_CONTEXT, http.MethodGet, TEST_BASE_URI, http.NoBody)
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Set(TEST_USER_ID_HEADER_KEY, TEST_USER_ID)
	expectedPayload, err := json.Marshal(trips)
	if err != nil {
		t.Error(err.Error())
		return
	}
	respWr.EXPECT().Write(expectedPayload).Return(len(expectedPayload), nil)
	respWr.EXPECT().WriteHeader(http.StatusOK)
	driverService.EXPECT().GetTrips(gomock.Any(), TEST_USER_ID).Return(trips)
	driverRouter.GetTripsHandler(respWr, req)
}

func TestWriteResponseBodyNotFullBodySentFailed(t *testing.T) {
	initDriverRouterMocks(t)
	req, err := http.NewRequestWithContext(TEST_BASE_CONTEXT, http.MethodGet, TEST_BASE_URI, http.NoBody)
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Set(TEST_USER_ID_HEADER_KEY, TEST_USER_ID)
	expectedPayload, err := json.Marshal(trips)
	if err != nil {
		t.Error(err.Error())
		return
	}
	respWr.EXPECT().Write(expectedPayload).Return(len(expectedPayload)-1, nil)
	respWr.EXPECT().WriteHeader(http.StatusInternalServerError)
	driverService.EXPECT().GetTrips(gomock.Any(), TEST_USER_ID).Return(trips)
	driverRouter.GetTripsHandler(respWr, req)
}

func TestWriteResponseBodyWriteErrFailed(t *testing.T) {
	initDriverRouterMocks(t)
	req, err := http.NewRequestWithContext(TEST_BASE_CONTEXT, http.MethodGet, TEST_BASE_URI, http.NoBody)
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Set(TEST_USER_ID_HEADER_KEY, TEST_USER_ID)
	expectedPayload, err := json.Marshal(trips)
	if err != nil {
		t.Error(err.Error())
		return
	}
	respWr.EXPECT().Write(expectedPayload).Return(0, errors.New("test_error"))
	respWr.EXPECT().WriteHeader(http.StatusInternalServerError)
	driverService.EXPECT().GetTrips(gomock.Any(), TEST_USER_ID).Return(trips)
	driverRouter.GetTripsHandler(respWr, req)
}

func TestGetTripHandlerSuccess(t *testing.T) {
	initDriverRouterMocks(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/trips/{trip_id}", driverRouter.GetTripHandler)
	req, err := http.NewRequestWithContext(TEST_BASE_CONTEXT, http.MethodGet, fmt.Sprintf("http://%s/trips/%s", TEST_BASE_URI, TEST_TRIP_ID), http.NoBody)
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Set(TEST_USER_ID_HEADER_KEY, TEST_USER_ID)
	expectedPayload, err := json.Marshal(trips[0])
	if err != nil {
		t.Error(err.Error())
		return
	}
	respWr.EXPECT().Write(expectedPayload).Return(len(expectedPayload), nil)
	respWr.EXPECT().WriteHeader(http.StatusOK)
	driverService.EXPECT().GetTrip(gomock.Any(), TEST_USER_ID, TEST_TRIP_ID).Return(&trips[0])
	rtr.ServeHTTP(respWr, req)
}

func TestGetTripHandlerNotFoundFailed(t *testing.T) {
	initDriverRouterMocks(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/trips/{trip_id}", driverRouter.GetTripHandler)
	req, err := http.NewRequestWithContext(TEST_BASE_CONTEXT, http.MethodGet, fmt.Sprintf("http://%s/trips/%s", TEST_BASE_URI, TEST_TRIP_ID), http.NoBody)
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Set(TEST_USER_ID_HEADER_KEY, TEST_USER_ID)
	respWr.EXPECT().WriteHeader(http.StatusNotFound)
	driverService.EXPECT().GetTrip(gomock.Any(), TEST_USER_ID, TEST_TRIP_ID).Return(nil)
	rtr.ServeHTTP(respWr, req)
}

type handler func(http.ResponseWriter, *http.Request)

func TestHandlersNoUserIdHeaderFailed(t *testing.T) {
	initDriverRouterMocks(t)
	req, err := http.NewRequest(http.MethodGet, TEST_BASE_URI, http.NoBody)
	if err != nil {
		t.Error(err.Error())
		return
	}
	handlers := []handler{
		driverRouter.GetTripsHandler,
		driverRouter.GetTripHandler,
		driverRouter.StartTripHandler,
		driverRouter.EndTripHandler,
		driverRouter.AcceptTripHandler,
		driverRouter.CancelTripHandler,
	}
	for _, handler := range handlers {
		respWr.EXPECT().WriteHeader(http.StatusBadRequest)
		handler(respWr, req)
	}
}

func TestCancelTripHandlerSuccess(t *testing.T) {
	initDriverRouterMocks(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/trips/{trip_id}/cancel", driverRouter.CancelTripHandler).Methods(http.MethodPost)
	req, err := http.NewRequestWithContext(TEST_BASE_CONTEXT, http.MethodPost,
		fmt.Sprintf("http://%s/trips/%s/cancel", TEST_BASE_URI, TEST_TRIP_ID), bytes.NewReader([]byte(TEST_CANCEL_REASON)))
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Set(TEST_USER_ID_HEADER_KEY, TEST_USER_ID)
	respWr.EXPECT().WriteHeader(http.StatusOK)
	driverService.EXPECT().CancelTrip(gomock.Any(), TEST_USER_ID, TEST_TRIP_ID, TEST_CANCEL_REASON).Return(svc.Ok)
	rtr.ServeHTTP(respWr, req)
}

func TestCancelTripHandlerTripNotFoundFailed(t *testing.T) {
	initDriverRouterMocks(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/trips/{trip_id}/cancel", driverRouter.CancelTripHandler).Methods(http.MethodPost)
	req, err := http.NewRequestWithContext(TEST_BASE_CONTEXT, http.MethodPost,
		fmt.Sprintf("http://%s/trips/%s/cancel", TEST_BASE_URI, TEST_TRIP_ID), bytes.NewReader([]byte(TEST_CANCEL_REASON)))
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Set(TEST_USER_ID_HEADER_KEY, TEST_USER_ID)
	respWr.EXPECT().WriteHeader(http.StatusNotFound)
	driverService.EXPECT().CancelTrip(gomock.Any(), TEST_USER_ID, TEST_TRIP_ID, TEST_CANCEL_REASON).Return(svc.TripNotFound)
	rtr.ServeHTTP(respWr, req)
}

func TestCancelTripHandlerInvalidOpFailed(t *testing.T) {
	initDriverRouterMocks(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/trips/{trip_id}/cancel", driverRouter.CancelTripHandler).Methods(http.MethodPost)
	req, err := http.NewRequestWithContext(TEST_BASE_CONTEXT, http.MethodPost,
		fmt.Sprintf("http://%s/trips/%s/cancel", TEST_BASE_URI, TEST_TRIP_ID), bytes.NewReader([]byte(TEST_CANCEL_REASON)))
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Set(TEST_USER_ID_HEADER_KEY, TEST_USER_ID)
	respWr.EXPECT().WriteHeader(http.StatusBadRequest)
	driverService.EXPECT().CancelTrip(gomock.Any(), TEST_USER_ID, TEST_TRIP_ID, TEST_CANCEL_REASON).Return(svc.InvalidNewStatus)
	rtr.ServeHTTP(respWr, req)
}

func TestCancelTripHandlerInternalErrorFailed(t *testing.T) {
	initDriverRouterMocks(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/trips/{trip_id}/cancel", driverRouter.CancelTripHandler).Methods(http.MethodPost)
	req, err := http.NewRequestWithContext(TEST_BASE_CONTEXT, http.MethodPost,
		fmt.Sprintf("http://%s/trips/%s/cancel", TEST_BASE_URI, TEST_TRIP_ID), bytes.NewReader([]byte(TEST_CANCEL_REASON)))
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Set(TEST_USER_ID_HEADER_KEY, TEST_USER_ID)
	respWr.EXPECT().WriteHeader(http.StatusInternalServerError)
	driverService.EXPECT().CancelTrip(gomock.Any(), TEST_USER_ID, TEST_TRIP_ID, TEST_CANCEL_REASON).Return(svc.InternalError)
	rtr.ServeHTTP(respWr, req)
}

type testTripCase struct {
	Handler         handler
	DriverSvcMethod *gomock.Call
	HandlerUriName  string
}

func initTestTripHandlerCases() []testTripCase {
	return []testTripCase{
		{
			Handler:         driverRouter.AcceptTripHandler,
			DriverSvcMethod: driverService.EXPECT().AcceptTrip(gomock.Any(), TEST_USER_ID, TEST_TRIP_ID),
			HandlerUriName:  "accept",
		},
		{
			Handler:         driverRouter.StartTripHandler,
			DriverSvcMethod: driverService.EXPECT().StartTrip(gomock.Any(), TEST_USER_ID, TEST_TRIP_ID),
			HandlerUriName:  "start",
		},
		{
			Handler:         driverRouter.EndTripHandler,
			DriverSvcMethod: driverService.EXPECT().EndTrip(gomock.Any(), TEST_USER_ID, TEST_TRIP_ID),
			HandlerUriName:  "end",
		},
	}
}

func TestTripHandlerSuccess(t *testing.T) {
	initDriverRouterMocks(t)
	simpleTripOpAction(t, svc.Ok, http.StatusOK)
}

func TestTripHandlerTripNotFoundFailed(t *testing.T) {
	initDriverRouterMocks(t)
	simpleTripOpAction(t, svc.TripNotFound, http.StatusNotFound)
}

func TestTripHandlerTripInvalidOpFailed(t *testing.T) {
	initDriverRouterMocks(t)
	simpleTripOpAction(t, svc.InvalidNewStatus, http.StatusBadRequest)
}

func TestTripHandlerTripInternalErrFailed(t *testing.T) {
	initDriverRouterMocks(t)
	simpleTripOpAction(t, svc.InternalError, http.StatusInternalServerError)
}

func simpleTripOpAction(t *testing.T, returnedChangeTripOpStatus svc.ChangeTripStatusStatus, expectedHttpStatus int) {
	cases := initTestTripHandlerCases()
	for _, cs := range cases {
		rtr := mux.NewRouter()
		rtr.HandleFunc("/trips/{trip_id}/"+cs.HandlerUriName, cs.Handler).Methods(http.MethodPost)
		req, err := http.NewRequestWithContext(TEST_BASE_CONTEXT, http.MethodPost,
			fmt.Sprintf("http://%s/trips/%s/%s", TEST_BASE_URI, TEST_TRIP_ID, cs.HandlerUriName), http.NoBody)
		if err != nil {
			t.Error(err.Error())
			return
		}
		req.Header.Set(TEST_USER_ID_HEADER_KEY, TEST_USER_ID)
		respWr.EXPECT().WriteHeader(expectedHttpStatus)
		cs.DriverSvcMethod.Return(returnedChangeTripOpStatus)
		rtr.ServeHTTP(respWr, req)
	}
}
