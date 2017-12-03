package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"

	"gopkg.in/mgo.v2/bson"
)

func FulfillmentHandler(w http.ResponseWriter, req *http.Request) {
	c := session.DB("dispotrains").C("stations")
	jsonReader := json.NewDecoder(req.Body)
	data := make(map[string]interface{})
	jsonReader.Decode(&data)
	queryResult, ok := data["queryResult"].(map[string]interface{})
	if !ok {
		log.Println(data["queryResult"])
		return
	}
	action, ok := queryResult["action"].(string)
	if !ok {
		log.Println(queryResult["action"])
		return
	}
	if action != "get_station_info" {
		log.Println(errors.New("Unknown action: " + action))
		return
	}
	parameters, ok := queryResult["parameters"].(map[string]interface{})
	if !ok {
		log.Println(queryResult["parameters"])
		return
	}
	stationName, ok := parameters["station"].(string)
	if !ok {
		log.Println(parameters["station"])
		return
	}
	var station storage.Station
	if err := c.Find(bson.M{"name": stationName}).One(&station); err != nil {
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(&station)
}
