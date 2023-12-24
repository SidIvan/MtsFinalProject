package app

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"

	"gitlab.com/AntYats/go_project/internal/httpadapter"
	// "github.com/anonimpopov/hw4/internal/service"
)

const (
	AppName                   = "TaxiService"
	DefaultServeAddress       = "localhost:8080"
	DefaultShutdownTimeout    = 20 * time.Second
	DefaultBasePath           = "/location/"
	DefaultAccessTokenCookie  = "access_token"
	DefaultRefreshTokenCookie = "refresh_token"
	DefaultDSN                = "dsn://"
	DefaultMigrationsDir      = "file://migrations/"
)

type AppConfig struct {
	Debug           bool          `yaml:"debug"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type Config struct {
	App      AppConfig          `yaml:"app"`
	Database DatabaseConfig     `yaml:"database"`
	HTTP     httpadapter.Config `yaml:"http"`
}

func NewConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cnf := Config{
		App: AppConfig{
			ShutdownTimeout: DefaultShutdownTimeout,
		},
		Database: DatabaseConfig{
			DSN:           DefaultDSN,
			MigrationsDir: DefaultMigrationsDir,
		},
		HTTP: httpadapter.Config{
			ServeAddress:       DefaultServeAddress,
			BasePath:           DefaultBasePath,
			AccessTokenCookie:  DefaultAccessTokenCookie,
			RefreshTokenCookie: DefaultRefreshTokenCookie,
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
