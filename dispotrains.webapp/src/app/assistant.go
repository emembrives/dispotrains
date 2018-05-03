package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"

	"gopkg.in/mgo.v2/bson"
)

var (
	months = map[time.Month]string{
		1:  "janvier",
		2:  "février",
		3:  "mars",
		4:  "avril",
		5:  "mai",
		6:  "juin",
		7:  "juillet",
		8:  "août",
		9:  "septembre",
		10: "octobre",
		11: "novembre",
		12: "décembre",
	}
)

type webhookResponse struct {
	FulfillmentText string `json:"fulfillmentText,omitempty"`
}

func brokenElevators(elevators []*storage.Elevator) []*storage.Elevator {
	broken := make([]*storage.Elevator, 0)
	for _, e := range elevators {
		state := e.GetLastStatus().State
		if state != "Disponible" && state != "Information non disponible" {
			broken = append(broken, e)
		}
	}
	return broken
}

func noInfoElevators(elevators []*storage.Elevator) []*storage.Elevator {
	noInfo := make([]*storage.Elevator, 0)
	for _, e := range elevators {
		state := e.GetLastStatus().State
		if state == "Information non disponible" {
			noInfo = append(noInfo, e)
		}
	}
	return noInfo
}

func FulfillmentHandler(w http.ResponseWriter, req *http.Request) {
	c := session.DB("dispotrains").C("stations")
	jsonReader := json.NewDecoder(req.Body)
	data := make(map[string]interface{})
	jsonReader.Decode(&data)
	queryResult, ok := data["queryResult"].(map[string]interface{})
	if !ok {
		log.Printf("Error decoding queryResult: %v\n", data["queryResult"])
		return
	}
	action, ok := queryResult["action"].(string)
	if !ok {
		log.Printf("Error decoding action: %v\n", queryResult["action"])
		return
	}
	if action != "get_station_info" {
		log.Println(errors.New("Unknown action: " + action))
		return
	}
	parameters, ok := queryResult["parameters"].(map[string]interface{})
	if !ok {
		log.Printf("Error decoding parameters: %v\n", queryResult["parameters"])
		return
	}
	stationName, ok := parameters["station"].(string)
	if !ok {
		log.Printf("Error decoding station: %v\n", parameters["station"])
		return
	}
	var station storage.Station
	if err := c.Find(bson.M{"name": stationName}).One(&station); err != nil {
		log.Printf("Station %s not found: %v", stationName, err)
		return
	}
	response := webhookResponse{}
	response.FulfillmentText = fmt.Sprintf(
		"Au %d %s à %s, la gare de %s",
		station.LastElevatorUpdate().Day(),
		months[station.LastElevatorUpdate().Month()],
		station.LastElevatorUpdate().Format("15 heures 04"),
		station.DisplayName)
	if station.Available() {
		response.FulfillmentText = response.FulfillmentText + " n'a aucun ascenseur en panne."
	} else {
		broken := brokenElevators(station.Elevators)
		noInfo := noInfoElevators(station.Elevators)
		var summary string
		var descriptions string
		if len(broken) == 0 && len(noInfo) == 0 {
			summary = fmt.Sprintf(" a %d ascenseurs sur %d sans relevé.", len(noInfo), len(station.Elevators))
		} else if len(broken) != 0 {
			summary = fmt.Sprintf(" a %d ascenseurs sur %d en panne et %d sans relevé. ", len(broken), len(station.Elevators), len(noInfo))
			descriptions = "Les ascenseurs en panne sont: "
			for i, e := range broken {
				if i != 0 {
					descriptions += " ; "
				}
				descriptions += fmt.Sprintf("%s", e.Situation)
				if len(e.Direction) != 0 {
					descriptions += ", direction " + e.Direction
				}
				if e.Status.Forecast != nil {
					descriptions += ", jusqu'au " + e.Status.Forecast.Format("02/01/2006")
				}
			}
		}
		response.FulfillmentText = response.FulfillmentText + summary + descriptions
	}
	json.NewEncoder(w).Encode(&response)
}
