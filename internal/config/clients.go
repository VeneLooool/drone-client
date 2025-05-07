package config

import (
	"context"
	"os"
)

const (
	EnvKeyDronesApiHost     = "DRONES_API_HOST"
	EnvKeyDronesApiGrpcPort = "DRONES_API_GRPC_PORT"
)

type DroneApiClientConfig struct {
	Host string
	Port string
}

func NewDroneApiClientConfig(ctx context.Context) *DroneApiClientConfig {
	return &DroneApiClientConfig{
		Host: os.Getenv(EnvKeyDronesApiHost),
		Port: os.Getenv(EnvKeyDronesApiGrpcPort),
	}
}
