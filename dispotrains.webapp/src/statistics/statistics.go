package statistics

import (
	"log"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type storageStatus struct {
	Elevator   string
	State      string
	Lastupdate time.Time
}

type NetworkStats struct {
	Good    int
	Bad     int
	LongBad int
}

func newElevatorState(status storageStatus) *storage.ElevatorState {
	state := &storage.ElevatorState{}
	state.Elevator = status.Elevator
	state.State = status.State
	state.Begin = status.Lastupdate
	state.End = status.Lastupdate
	return state
}

func toUpsertQuery(es *storage.ElevatorState) bson.M {
	return bson.M{
		"elevator": es.Elevator,
		"state":    es.State,
		"begin":    es.Begin}
}

// ComputeElevatorStatistics computes and stores per-elevator statistics.
func ComputeElevatorStatistics(session *mgo.Session) error {
	cStatuses := session.DB("dispotrains").C("statuses")
	err := cStatuses.EnsureIndexKey("elevator", "lastupdate")
	if err != nil {
		return err
	}

	var elevators []string
	err = cStatuses.Find(nil).Select(bson.M{"elevator": 1}).Distinct("elevator", &elevators)
	if err != nil {
		return err
	}

	cStatistics := session.DB("dispotrains").C("statistics")

	for _, elevatorID := range elevators {
		log.Printf("Processing elevator %s\n", elevatorID)
		elevatorState := &storage.ElevatorState{}
		err = cStatistics.Find(bson.M{"elevator": elevatorID}).Sort("-end").Limit(1).One(elevatorState)
		if err != nil && err != mgo.ErrNotFound {
			return err
		} else if err == mgo.ErrNotFound {
			elevatorState = nil
		}
		query := bson.M{"elevator": elevatorID}
		if elevatorState != nil {
			query["lastupdate"] = bson.M{"$gt": elevatorState.End}
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
			if _, isUnknown := storage.UnknownStates[status.State]; isUnknown {
				continue
			}
			if (status.State == "Disponible") != (elevatorState.State == "Disponible") {
				cStatistics.Upsert(toUpsertQuery(elevatorState), elevatorState)
				elevatorState = newElevatorState(status)
			}
		}
		cStatistics.Upsert(toUpsertQuery(elevatorState), elevatorState)
		if err = iter.Close(); err != nil {
			return err
		}
	}

	err = cStatistics.EnsureIndexKey("elevator", "end")
	return err
}

func ComputeGlobalStatistics(session *mgo.Session) (*NetworkStats, error) {
	cStatistics := session.DB("dispotrains").C("statistics")
	results := make([]bson.M, 0)
	pipe := cStatistics.Pipe(
		[]bson.M{
			bson.M{"$match": bson.M{"end": bson.M{"$gte": time.Now().AddDate(0, 0, -10)}}},
			bson.M{"$sort": bson.M{"end": 1}},
			bson.M{"$group": bson.M{
				"_id":   "$elevator",
				"state": bson.M{"$last": "$state"},
				"begin": bson.M{"$last": "$begin"},
				"end":   bson.M{"$last": "$end"}}},
		},
	)
	err := pipe.All(&results)
	if err != nil {
		return nil, err
	}
	ns := NetworkStats{}
	longLimit := time.Now().AddDate(0, 0, -3)
	for _, elevator := range results {
		if elevator["state"] == "Disponible" {
			ns.Good++
		} else {
			beginTime, ok := elevator["begin"].(time.Time)
			if !ok {
				panic(elevator["begin"])
			}
			ns.Bad++
			if beginTime.Before(longLimit) {
				ns.LongBad++
			}
		}
	}
	return &ns, nil
}
