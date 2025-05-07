package drone

import (
	"context"
	"log"
	"time"
	
	"github.com/VeneLooool/drone-client/internal/model"
)

type UseCase struct {
	publisher Publisher
}

func New(publisher Publisher) *UseCase {
	return &UseCase{
		publisher: publisher,
	}
}

func (u *UseCase) StartDroneMission(ctx context.Context, droneID uint64, mission model.Mission) (err error) {
	go func() {
		u.testMission(context.Background(), droneID, mission)
	}()
	return nil
}

func (u *UseCase) testMission(ctx context.Context, droneID uint64, mission model.Mission) {
	testEvents := []model.Event{
		{
			Drone:     model.Drone{ID: droneID, Status: model.DroneStatusInMission},
			EventType: model.EventTypeDroneChangeStatus,
		},
		{
			Drone:     model.Drone{ID: droneID, Status: model.DroneStatusCharging},
			EventType: model.EventTypeDroneChangeStatus,
		},
		{
			Drone:     model.Drone{ID: droneID, Status: model.DroneStatusAvailable},
			EventType: model.EventTypeDroneChangeStatus,
		},
	}

	for _, event := range testEvents {
		if err := u.publisher.Publish(ctx, event); err != nil {
			log.Printf("Failed to publish event, error: %s", err.Error())
		}
		time.Sleep(30 * time.Second)
	}
}
