package httpadapter

import (
	"context"
	// "fmt"
	"encoding/json"
	"net/http"

	// httpSwagger "github.com/swaggo/http-swagger"

	// "gitlab.com/AntYats/go_project/internal/docs" // go:generate
	"github.com/go-chi/chi/v5"
	"gitlab.com/AntYats/go_project/internal/model"
	"gitlab.com/AntYats/go_project/internal/service"
	"go.uber.org/zap"
)

type adapter struct {
	config *Config

	locationService service.User

	server *http.Server
	logger *zap.Logger
}

func (a *adapter) GetDrivers(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")

	if lat == "" {
		a.logger.Error("Lat not found")
		return
	}

	lng := r.URL.Query().Get("lng")

	if lng == "" {
		a.logger.Error("Lng not found")
		writeError(w, ErrNotFound)
		return
	}

	radius_string := r.URL.Query().Get("radius")

	if radius_string == "" {
		a.logger.Error("Radius not found")
		writeError(w, ErrNotFound)
		return
	}

	coords := model.UserData{
		Lat: convertToFloat(lat),
		Lng: convertToFloat(lng),
	}

	radius := convertToFloat(radius_string)

	a.logger.Info("Getting drivers")

	users, err := a.locationService.GetDrivers(r.Context(), float64(radius), &coords)

	if err != nil {
		a.logger.Error("Failed to get drivers nearby")
		writeError(w, err)
		return
	}

	writeJSONResponse(w, http.StatusAccepted, users)

}

func (a *adapter) ChangeDriverInfo(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Starting to change driver info")

	var coords model.UserData

	a.logger.Info("Parsing current driver location")
	parse_err := json.NewDecoder(r.Body).Decode(&coords)

	if parse_err != nil {
		a.logger.Error("Failed to parse current driver location")
		writeError(w, parse_err)
		return
	}

	user_id := chi.URLParam(r, "user_id")

	err := a.locationService.ChangeDriverInfo(r.Context(), user_id, &coords)

	if err != nil {
		a.logger.Error("Failed to change current driver info")
		writeError(w, err)
		return
	}

	a.logger.Info("Successfully changed driver info")

	writeJSONResponse(w, http.StatusAccepted, "DONE")

}

func (a *adapter) Serve() error {
	a.logger.Info("Serving started")

	r := chi.NewRouter()

	apiRouter := chi.NewRouter()
	apiRouter.Post("/drivers/{user_id}/location", a.ChangeDriverInfo)
	apiRouter.Get("/drivers", a.GetDrivers)

	r.Mount(a.config.BasePath, apiRouter)

	a.server = &http.Server{Addr: a.config.ServeAddress, Handler: r}

	if a.config.UseTLS {
		return a.server.ListenAndServeTLS(a.config.TLSCrtFile, a.config.TLSKeyFile)
	}

	return a.server.ListenAndServe()
}

func (a *adapter) Shutdown(ctx context.Context) {
	_ = a.server.Shutdown(ctx)
}

func New(config *Config, l service.User) Adapter {
	// if config.SwaggerAddress != "" {
	// 	docs.SwaggerInfo.Host = config.SwaggerAddress
	// } else {
	// 	docs.SwaggerInfo.Host = config.ServeAddress
	// }

	// docs.SwaggerInfo.BasePath = config.BasePath

	return &adapter{
		config:          config,
		locationService: l,
	}
}
