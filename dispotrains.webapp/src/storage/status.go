package storage

import (
	"time"
)

type Status struct {
	State      string     `json:"state"`
	LastUpdate time.Time  `json:"lastupdate"`
	Forecast   *time.Time `json:"forecast"`
	elevator   *Elevator
}

func (s *Status) ToStorage() *StatusStorage {
	return &StatusStorage{
		ElevatorID: s.elevator.ID,
		LastUpdate: s.LastUpdate,
		State:      s.State,
		Forecast:   s.Forecast,
	}
}

type StatusStorage struct {
	ElevatorID string     `json:"elevatorid"`
	LastUpdate time.Time  `json:"lastupdate"`
	State      string     `json:"state"`
	Forecast   *time.Time `json:"forecast,omitempty"`
}
