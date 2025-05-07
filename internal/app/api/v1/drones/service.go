package drones

import (
	desc "github.com/VeneLooool/drone-client/internal/pb/api/v1/drones"
)

// Implementation is a Service implementation
type Implementation struct {
	desc.UnimplementedDronesServer
	DroneUC DroneUC
}

// NewService return new instance of Implementation.
func NewService(DroneUC DroneUC) *Implementation {
	return &Implementation{
		DroneUC: DroneUC,
	}
}
