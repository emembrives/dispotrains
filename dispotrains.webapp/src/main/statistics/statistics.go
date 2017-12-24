package main

import (
	"log"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	server = "db"
)

type storageStatus struct {
	Elevator   string
	State      string
	Lastupdate time.Time
}

func newElevatorState(status storageStatus) *storage.ElevatorState {
	state := &storage.ElevatorState{}
	state.Elevator = status.Elevator
	state.State = status.State
	state.Begin = status.Lastupdate
	state.End = status.Lastupdate
	return state
}

func main() {
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cStatuses := session.DB("dispotrains").C("statuses")
	err = cStatuses.EnsureIndexKey("lastupdate")
	if err != nil {
		panic(err)
	}
	err = cStatuses.EnsureIndexKey("elevator")
	if err != nil {
		panic(err)
	}

	var elevators []string
	err = cStatuses.Find(nil).Select(bson.M{"elevator": 1}).Distinct("elevator", &elevators)
	if err != nil {
		panic(err)
	}

	cStatistics := session.DB("dispotrains").C("statistics")

	for _, elevatorID := range elevators {
		log.Printf("Processing elevator %s\n", elevatorID)
		elevatorState := &storage.ElevatorState{}
		err = cStatistics.Find(bson.M{"elevator": elevatorID}).Sort("-begin").Limit(1).One(elevatorState)
		if err != nil && err != mgo.ErrNotFound {
			panic(err)
		} else if err == mgo.ErrNotFound {
			elevatorState = nil
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
			panic(err)
		}
	}

	err = cStatistics.EnsureIndexKey("elevator")
	if err != nil {
		panic(err)
	}
}
