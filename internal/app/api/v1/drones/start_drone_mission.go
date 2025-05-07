package drones

import (
	"context"

	"github.com/VeneLooool/drone-client/internal/model"
	desc "github.com/VeneLooool/drone-client/internal/pb/api/v1/drones"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) StartDroneMission(ctx context.Context, req *desc.StartDroneMission_Request) (*empty.Empty, error) {
	if err := i.DroneUC.StartDroneMission(ctx, req.GetId(), transformMissionToModel(req.GetMission())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func transformMissionToModel(mission *desc.Mission) model.Mission {
	return model.Mission{
		Coordinates: transformCoordinatesToModel(mission.GetCoordinates()),
	}
}

func transformCoordinatesToModel(protoCoordinates []*desc.Coordinate) model.Coordinates {
	coordinates := make([]model.Coordinate, 0, len(protoCoordinates))
	for _, coordinate := range protoCoordinates {
		coordinates = append(coordinates, model.Coordinate{
			Latitude:  coordinate.Latitude,
			Longitude: coordinate.Longitude,
		})
	}
	return coordinates
}
