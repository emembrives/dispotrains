package storage

import (
	"time"
)

type Status struct {
	State      string
	LastUpdate time.Time
	Forecast   *time.Time
	elevator   *Elevator
}
