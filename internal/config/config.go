package config

import (
	"context"
	"os"
)

const (
	EnvKeyHttpPort = "DRONE_CLIENT_HTTP_PORT"
	EnvKeyGrpcPort = "DRONE_CLIENT_GRPC_PORT"

	EnvKeyKafkaHost = "KAFKA_INTERNAL_HOST"
	EnvKeyKafkaPort = "KAFKA_INTERNAL_PORT"
)

type Config struct {
	HttpPort string
	GrpcPort string

	DroneApiClientConfig *DroneApiClientConfig
	KafkaConfig          *KafkaConfig
}

type KafkaConfig struct {
	KafkaHost string
	KafkaPort string
}

func (c *Config) GetKafkaConfig() *KafkaConfig {
	return c.KafkaConfig
}

func New(ctx context.Context) (*Config, error) {
	return &Config{
		HttpPort: os.Getenv(EnvKeyHttpPort),
		GrpcPort: os.Getenv(EnvKeyGrpcPort),

		DroneApiClientConfig: NewDroneApiClientConfig(ctx),

		KafkaConfig: &KafkaConfig{
			KafkaPort: os.Getenv(EnvKeyKafkaPort),
			KafkaHost: os.Getenv(EnvKeyKafkaHost),
		},
	}, nil
}
