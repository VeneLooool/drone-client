package model

type DroneStatus string

func (ds DroneStatus) String() string {
	return string(ds)
}

const (
	DroneStatusAvailable   DroneStatus = "available"
	DroneStatusInMission   DroneStatus = "mission"
	DroneStatusCharging    DroneStatus = "charging"
	DroneStatusMaintenance DroneStatus = "maintenance"
	DroneStatusOffline     DroneStatus = "offline"
)

type Drone struct {
	ID     uint64      `db:"id" json:"id"`
	Status DroneStatus `db:"status" json:"status"`
}

func (d *Drone) SetDefaultStatus() {
	d.Status = DroneStatusAvailable
}

func (d *Drone) SetStatus(status DroneStatus) {
	d.Status = status
}
