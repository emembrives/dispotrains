package storage

import (
	"strings"
	"time"
)

type Elevator struct {
	ID        string
	Situation string
	Direction string
	station   *Station
	Status    *Status
}

func (elevator *Elevator) GetStation() *Station {
	return elevator.station
}

func (elevator *Elevator) NewStatus(description string, date string) (*Status, error) {
	var status *Status = &Status{State: strings.TrimSpace(description),
		elevator: elevator}
	lastUpdate, err := time.Parse("02/01/2006 15:04", strings.TrimSpace(date))
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		panic(err)
	}
	status.LastUpdate = overrideLocation(lastUpdate, loc)
	elevator.Status = status
	return status, nil
}

func (elevator *Elevator) GetLastStatus() *Status {
	return elevator.Status
}

func overrideLocation(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.Date()
	H, M, S := t.Clock()
	return time.Date(y, m, d, H, M, S, t.Nanosecond(), loc)
}
