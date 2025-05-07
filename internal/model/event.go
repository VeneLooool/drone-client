package model

import "strconv"

type EventType string

func (et EventType) String() string {
	return string(et)
}

const (
	EventTypeDroneChangeStatus EventType = "drone_change_status"
)

type Event struct {
	EventType EventType `json:"event_type"`
	Drone     Drone     `json:"drone"`
}

func (e Event) GetEventKey() string {
	return strconv.FormatUint(e.Drone.ID, 10)
}
