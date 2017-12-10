package main

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	server = "localhost"
)

type elevatorState struct {
	Elevator string
	State    string
	Begin    time.Time
	End      time.Time
}

type storageStatus struct {
	Elevator   string
	State      string
	Lastupdate time.Time
}

func newElevatorState(status storageStatus) *elevatorState {
	state := &elevatorState{}
	state.Elevator = status.Elevator
	state.State = status.State
	state.Begin = status.Lastupdate
	state.End = status.Lastupdate
	return state
}

func main() {
	session, err := mgo.Dial(server)
	if err != nil {
		log.Fatalln(err)
	}
	defer session.Close()

	cStatuses := session.DB("dispotrains").C("statuses")
	err = cStatuses.EnsureIndexKey("lastupdate")
	if err != nil {
		log.Fatalln(err)
	}
	err = cStatuses.EnsureIndexKey("elevator")
	if err != nil {
		log.Fatalln(err)
	}

	var elevators []string
	err = cStatuses.Find(nil).Select(bson.M{"elevator": 1}).Distinct("elevator", &elevators)
	if err != nil {
		log.Fatalln(err)
	}

	cStatistics := session.DB("dispotrains").C("statistics")

	for _, elevatorID := range elevators {
		log.Printf("Processing elevator %s\n", elevatorID)
		var elevatorState *elevatorState
		err = cStatistics.Find(bson.M{"elevator": elevatorID}).Sort("-begin").Limit(1).One(elevatorState)
		if err != nil && err != mgo.ErrNotFound {
			log.Fatalln(err)
		}
		query := bson.M{"elevator": elevatorID}
		if elevatorState != nil {
			query["lastupdate"] = bson.M{"$gte": elevatorState.Begin}
		}
		iter := cStatuses.Find(query).Sort("lastupdate").Iter()
		var status storageStatus
		for iter.Next(&status) {
			if elevatorState == nil {
				elevatorState = newElevatorState(status)
				continue
			}
			if status.Lastupdate.Before(elevatorState.Begin) {
				continue
			}
			elevatorState.End = status.Lastupdate
			if status.State != elevatorState.State {
				cStatistics.Insert(elevatorState)
				elevatorState = newElevatorState(status)
			}
		}
		if err := iter.Close(); err != nil {
			log.Fatalln(err)
		}
	}

	err = cStatistics.EnsureIndexKey("elevator")
	if err != nil {
		log.Fatalln(err)
	}
}
