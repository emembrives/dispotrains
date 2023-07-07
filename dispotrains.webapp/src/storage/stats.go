package storage

import "time"

type ElevatorState struct {
	Elevator string    `json:"elevator"`
	State    string    `json:"state"`
	Begin    time.Time `json:"begin"`
	End      time.Time `json:"end"`
}
