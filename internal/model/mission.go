package model

type Coordinates []Coordinate

type Coordinate struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Mission struct {
	Coordinates Coordinates `db:"coordinates"`
}
