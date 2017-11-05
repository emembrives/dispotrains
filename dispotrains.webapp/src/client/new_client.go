package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
)

type parser struct {
	stations map[string]*storage.Station
}

func NewGetAllLines() ([]*storage.Line, error) {
	req, err := http.NewRequest("GET",
		"https://api.vianavigo.com/elevatorsInfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Host-Override", "vgo-api")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	parser := newParser()
	return parser.parseRawData(resp.Body)
}

func newParser() *parser {
	parser := &parser{}
	parser.stations = make(map[string]*storage.Station)
	return parser
}

func (parser *parser) parseRawData(input io.Reader) ([]*storage.Line, error) {
	decoder := json.NewDecoder(input)
	data := make([]map[string]interface{}, 0)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	lines := make([]*storage.Line, len(data))
	var err error
	for index, lineData := range data {
		lines[index], err = parser.parseLine(lineData)
		if err != nil {
			return nil, err
		}
	}
	return lines, nil
}

func (parser *parser) parseLine(lineData map[string]interface{}) (*storage.Line, error) {
	network, ok := lineData["mode"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"lineData[\"mode\"]=%+v not a string", lineData["mode"])
	}
	name, ok := lineData["shortName"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"lineData[\"shortName\"]=%+v not a string", lineData["shortName"])
	}
	code, ok := lineData["id"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"lineData[\"code\"]=%+v not a string", lineData["code"])
	}
	line := storage.NewLine(network, name, code)
	stationsData, ok := lineData["stops"].([]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"Unable to interpret stations data for line %s", line.ID)
	}
	stations := make([]*storage.Station, len(stationsData))
	var err error
	for index, iStationData := range stationsData {
		stationData, ok := iStationData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(
				"Unable to parse station data for line %s", line.ID)
		}
		stations[index], err = parser.parseStation(stationData)
		if err != nil {
			return nil, err
		}
		stations[index].AttachLine(line)
	}
	return line, nil
}

func (parser *parser) parseStation(stationData map[string]interface{}) (*storage.Station, error) {
	name, ok := stationData["label"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"stationData[\"label\"]=%+v not a string", stationData["label"])
	}
	code, ok := stationData["id"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"stationData[\"id\"]=%+v not a string for station %s", stationData["id"],
			name)
	}
	var station *storage.Station
	if station, ok = parser.stations[code]; ok {
		return station, nil
	}
	station = storage.NewStation(name, "", code)
	parser.stations[code] = station
	elevatorsData, ok := stationData["elevators"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to parse elevators data")
	}
	var err error
	for _, iElevatorData := range elevatorsData {
		elevatorData, ok := iElevatorData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unable to parse elevator data")
		}
		_, err = parser.parseElevator(station, elevatorData)
		if err != nil {
			return nil, err
		}
	}
	return station, nil
}

func (parser *parser) parseElevator(
	station *storage.Station, elevatorData map[string]interface{}) (*storage.Elevator, error) {
	var ok bool
	elevatorID, ok := elevatorData["label"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"elevatorData[\"label\"]=%+v not a string", elevatorData["label"])
	}
	var elevatorSituation string
	if _, ok = elevatorData["situation"]; ok {
		elevatorSituation, ok = elevatorData["situation"].(string)
		if !ok {
			return nil, fmt.Errorf(
				"elevatorData[\"situation\"]=%+v not a string for elevator %s",
				elevatorData["situation"], elevatorID)
		}
	}
	var elevatorDirection string
	if _, ok = elevatorData["direction"]; ok {
		elevatorDirection, ok = elevatorData["direction"].(string)
		if !ok {
			return nil, fmt.Errorf(
				"elevatorData[\"direction\"]=%+v not a string for elevator %s",
				elevatorData["direction"], elevatorID)
		}
	}

	elevator := station.NewElevator(elevatorID, elevatorSituation,
		elevatorDirection)

	// Status
	if _, ok := elevatorData["stateUpdate"]; !ok {
		return elevator, nil
	}

	stateUpdate, ok := elevatorData["stateUpdate"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"elevatorData[\"stateUpdate\"]=%+v not a string for elevator %s", elevatorData["stateUpdate"],
			elevator.ID)
	}

	if _, ok := elevatorData["state"]; !ok {
		return elevator, nil
	}
	state, ok := elevatorData["state"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"elevatorData[\"state\"]=%+v not a string for elevator %s",
			elevatorData["state"], elevator.ID)
	}
	forecastStr := ""
	if _, ok = elevatorData["stateEndPrevision"]; ok {
		forecastStr, ok = elevatorData["stateEndPrevision"].(string)
		if !ok {
			return nil, fmt.Errorf(
				"elevatorData[\"stateEndPrevision\"]=%+v not a string for elevator %s",
				elevatorData["stateEndPrevision"], elevator.ID)
		}
	}
	elevator.NewViaNavigoStatus(state, stateUpdate, forecastStr)
	return elevator, nil
}
