package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/statistics"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
	"github.com/gorilla/mux"
	"github.com/linxGnu/grocksdb"
	"github.com/vmihailenco/msgpack/v5"
)

type outStats struct {
	Mtbf   time.Duration
	Mtbr   time.Duration
	Broken time.Duration
	Total  time.Duration
	States []*storage.ElevatorState
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// ElevatorHandle prepares the Elevator details page.
func ElevatorHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(req)
	elevatorID := vars["id"]

	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()

	iterator := session.NewIterator(ro)
	defer iterator.Close()

	stats := make([]*storage.ElevatorState, 0)
	datetime := time.Now().AddDate(-2, 0, 0)

	for iterator.SeekForPrev(storage.LastKeyOfBucket(storage.BucketStatistics, elevatorID)); iterator.ValidForPrefix(storage.MakeKey(storage.BucketStatistics, elevatorID)); iterator.Prev() {
		var state storage.ElevatorState
		err := msgpack.Unmarshal(iterator.Value().Data(), &state)
		if err != nil {
			log.Panic(err)
		}
		if state.End.Before(datetime) {
			break
		}
		stats = append(stats, &state)
	}

	var availableTime, brokenTime, totalTime time.Duration
	var availablePeriods, brokenPeriods int
	for _, stat := range stats {
		if stat.State == "Disponible" {
			availableTime += stat.End.Sub(stat.Begin)
			availablePeriods++
		} else {
			brokenTime += stat.End.Sub(stat.Begin)
			brokenPeriods++
		}
		totalTime += stat.End.Sub(stat.Begin)
	}
	out := outStats{}
	if availablePeriods != 0 {
		out.Mtbf = availableTime / time.Duration(availablePeriods)
	} else {
		out.Mtbf = 0
	}
	if brokenPeriods != 0 {
		out.Mtbr = brokenTime / time.Duration(brokenPeriods)
	} else {
		out.Mtbr = 0
	}
	out.Broken = brokenTime
	out.Total = totalTime
	out.States = stats[0:min(len(stats), 30)]
	if err := json.NewEncoder(w).Encode(&out); err != nil {
		log.Panic(err)
	}
}

func NetworkStatsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	out, err := statistics.ComputeGlobalStatistics(session)
	if err != nil {
		log.Panic(err)
	}
	if err := json.NewEncoder(w).Encode(&out); err != nil {
		log.Panic(err)
	}
}
