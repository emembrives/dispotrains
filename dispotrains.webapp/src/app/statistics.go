package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2/bson"
)

// ElevatorHandle prepares the Elevator details page.
func ElevatorHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(req)
	elevatorID := vars["id"]

	cStatistics := session.DB("dispotrains").C("statistics")

	stats := make([]storage.ElevatorState, 0)
	datetime := time.Now().AddDate(0, -1, 0)
	if err := cStatistics.Find(bson.M{"elevator": elevatorID, "enddate": bson.M{"$gt": datetime}}).All(&stats); err != nil {
		log.Println(err)
	}
	if err := json.NewEncoder(w).Encode(&stats); err != nil {
		log.Println(err)
	}

}
