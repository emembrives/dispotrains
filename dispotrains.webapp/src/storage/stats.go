package storage

import "time"

type ElevatorState struct {
	Elevator string
	State    string
	Begin    time.Time
	End      time.Time
}
