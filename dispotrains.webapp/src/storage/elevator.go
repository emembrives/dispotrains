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

	// UnknownStates lists unknown elevator states that may be ignored.
	UnknownStates = map[string]bool{
		"Inconnu ce jour":            true,
		"Information non disponible": true}
)

type Elevator struct {
	ID        string `json:"id"`
	Situation string `json:"situation"`
	Direction string `json:"direction"`
	station   *Station
	Status    *Status `json:"status"`
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
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		panic(err)
	}
	status.LastUpdate, err = time.ParseInLocation("2006-01-02T15:04", strings.TrimSpace(updateStr), loc)
	if err != nil {
		return nil, err
	}
	elevator.Status = status
	if len(forecastStr) != 0 {
		forecast, err := time.ParseInLocation("2006-01-02T15:04", strings.TrimSpace(forecastStr), loc)
		if err != nil {
			return nil, err
		}
		status.Forecast = &forecast
	} else if strings.Contains(state, " jusqu'au ") {
		forecastStr := state[len(state)-10:]
		t, err := time.Parse("02/01/2006", forecastStr)
		if err == nil {
			status.Forecast = &t
		}
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
