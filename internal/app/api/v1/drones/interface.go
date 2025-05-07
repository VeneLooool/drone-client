package drones

import (
	"context"

	"github.com/VeneLooool/drone-client/internal/model"
)

type DroneUC interface {
	StartDroneMission(ctx context.Context, droneID uint64, mission model.Mission) (err error)
}
