package config

import (
	"os"
	"strconv"
)

type LocationServiceConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	TimeoutSec int    `yaml:"timeout.sec"`
}

type MongoConfig struct {
	TimeoutSec                    int    `yaml:"timeout.sec"`
	URI                           string `yaml:"uri"`
	DatabaseName                  string `yaml:"database.name"`
	TripCollectionName            string `yaml:"trip.collection.name"`
	CancelReasonLogCollectionName string `yaml:"cancel.reason.log.collection.name"`
}

type KafkaConfig struct {
	Host              string `yaml:"host"`
	Port              string `yaml:"port"`
	ReadTopic         string `yaml:"read.topic"`
	WriteTopic        string `yaml:"write.topic"`
	GroupId           string `yaml:"group.id"`
	SessionTimeoutSec int    `yaml:"session.timeout.sec"`
	AsyncWrite        bool   `yaml:"is.async.write"`
	WriteBatchSize    int    `yaml:"write.batch.size"`
	TimeoutSec        int    `yaml:"timeout.sec"`
}

type DriverServiceConfig struct {
	EnvType                    string                 `yaml:"env"`
	Port                       int                    `yaml:"port"`
	Source                     string                 `yaml:"source"`
	GracefulShutdownTimeoutSec int                    `yaml:"timeout.sec"`
	SearchRadius               float64                `yaml:"search.radius"`
	LocationServiceConfig      *LocationServiceConfig `yaml:"location.service.config"`
	KafkaConfig                *KafkaConfig           `yaml:"kafka.config"`
	MongoConfig                *MongoConfig           `yaml:"mongo.config"`
	NatConfig                  *NatNotificationerConfig
}

func NewMongoConfigFromEnv() *MongoConfig {
	timeout, err := strconv.Atoi(os.Getenv("MONGO_CONFIG_TIMEOUT_SEC"))
	if err != nil {
		panic(err)
	}
	return &MongoConfig{
		TimeoutSec:                    timeout,
		URI:                           os.Getenv("MONGO_CONFIG_URI"),
		DatabaseName:                  os.Getenv("MONGO_CONFIG_DATABASE_NAME"),
		TripCollectionName:            os.Getenv("MONGO_CONFIG_TRIP_COLLECTION_NAME"),
		CancelReasonLogCollectionName: os.Getenv("MONGO_CONFIG_CANCEL_REASON_LOG_COLLECTION_NAME"),
	}
}

func NewKafkaConfigFromEnv() *KafkaConfig {
	sessionTimeout, err := strconv.Atoi(os.Getenv("KAFKA_CONFIG_SESSION_TIMEOUT_SEC"))
	if err != nil {
		panic(err)
	}
	batchSize, err := strconv.Atoi(os.Getenv("KAFKA_CONFIG_WRITE_BATCH_SIZE"))
	if err != nil {
		panic(err)
	}
	isAsync, err := strconv.ParseBool(os.Getenv("KAFKA_CONFIG_IS_ASYNC_WRITE"))
	if err != nil {
		panic(err)
	}
	timeout, err := strconv.Atoi(os.Getenv("KAFKA_CONFIG_TIMEOUT_SEC"))
	if err != nil {
		panic(err)
	}

	return &KafkaConfig{
		Host:              os.Getenv("KAFKA_CONFIG_HOST"),
		Port:              os.Getenv("KAFKA_CONFIG_PORT"),
		ReadTopic:         os.Getenv("KAFKA_CONFIG_READ_TOPIC"),
		WriteTopic:        os.Getenv("KAFKA_CONFIG_WRITE_TOPIC"),
		GroupId:           os.Getenv("KAFKA_CONFIG_GROUP_ID"),
		SessionTimeoutSec: sessionTimeout,
		AsyncWrite:        isAsync,
		WriteBatchSize:    batchSize,
		TimeoutSec:        timeout,
	}
}

func NewLocationSvcConfigFromEnv() *LocationServiceConfig {
	port, err := strconv.Atoi(os.Getenv("LOCATION_SERVICE_CONFIG_PORT"))
	if err != nil {
		panic(err)
	}
	timeout, err := strconv.Atoi(os.Getenv("LOCATION_SERVICE_CONFIG_TIMEOUT_SEC"))
	if err != nil {
		panic(err)
	}
	return &LocationServiceConfig{
		Host:       os.Getenv("LOCATION_SERVICE_CONFIG_HOST"),
		Port:       port,
		TimeoutSec: timeout,
	}
}

func NewDriverServiceConfigFromEnv() *DriverServiceConfig {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
	radius, err := strconv.ParseFloat(os.Getenv("RADIUS"), 64)
	if err != nil {
		panic(err)
	}
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT_SEC"))
	if err != nil {
		panic(err)
	}
	return &DriverServiceConfig{
		EnvType:                    os.Getenv("ENV"),
		Port:                       port,
		Source:                     os.Getenv("SOURCE"),
		SearchRadius:               radius,
		GracefulShutdownTimeoutSec: timeout,
		LocationServiceConfig:      NewLocationSvcConfigFromEnv(),
		KafkaConfig:                NewKafkaConfigFromEnv(),
		MongoConfig:                NewMongoConfigFromEnv(),
		NatConfig:                  NewNatNotificationerConfigFromEnv(),
	}
}

type NatNotificationerConfig struct {
	URI string
}

func NewNatNotificationerConfigFromEnv() *NatNotificationerConfig {
	return &NatNotificationerConfig{
		URI: os.Getenv("NAT_URI"),
	}
}
