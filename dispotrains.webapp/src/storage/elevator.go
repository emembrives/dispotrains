package storage

import (
	"strings"
	"time"
)

var (
	viaNavigoStates = map[string]string{
		"0": "Inconnu ce jour",
		"1": "Disponible",
		"2": "Hors-service",
		"3": "En travaux",
		"4": "Perturbation"}
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

func (elevator *Elevator) NewViaNavigoStatus(state string, updateStr string, forecastStr string) (*Status, error) {
	var status *Status = &Status{State: viaNavigoStates[state],
		elevator: elevator}
	lastUpdate, err := time.Parse("2006-01-02T15:04", strings.TrimSpace(updateStr))
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		panic(err)
	}
	status.LastUpdate = overrideLocation(lastUpdate, loc)
	elevator.Status = status
	if len(forecastStr) != 0 {
		forecast, err := time.Parse("2006-01-02T15:04", strings.TrimSpace(forecastStr))
		if err != nil {
			return nil, err
		}
		forecast = overrideLocation(forecast, loc)
		status.Forecast = &forecast
	}
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
