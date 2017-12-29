package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2/bson"
)

type outStats struct {
	Mtbf   time.Duration
	Mtbr   time.Duration
	Broken time.Duration
	Total  time.Duration
	States []storage.ElevatorState
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

	cStatistics := session.DB("dispotrains").C("statistics")

	stats := make([]storage.ElevatorState, 0)
	datetime := time.Now().AddDate(-2, 0, 0)
	if err := cStatistics.Find(bson.M{"elevator": elevatorID, "end": bson.M{"$gt": datetime}}).Sort("-begin").All(&stats); err != nil {
		log.Println(err)
	}
	availableStats := make([]storage.ElevatorState, 0)
	brokenStats := make([]storage.ElevatorState, 0)
	var broken, total time.Duration
	for _, stat := range stats {
		if stat.State == "Disponible" {
			availableStats = append(availableStats, stat)
		} else {
			brokenStats = append(brokenStats, stat)
			broken += stat.End.Sub(stat.Begin)
		}
		total += stat.End.Sub(stat.Begin)
	}
	sort.Slice(availableStats, func(i, j int) bool {
		iLen := availableStats[i].End.Sub(availableStats[i].Begin)
		jLen := availableStats[j].End.Sub(availableStats[j].Begin)
		return iLen < jLen
	})
	sort.Slice(brokenStats, func(i, j int) bool {
		iLen := brokenStats[i].End.Sub(brokenStats[i].Begin)
		jLen := brokenStats[j].End.Sub(brokenStats[j].Begin)
		return iLen < jLen
	})
	out := outStats{}
	if len(availableStats) != 0 {
		out.Mtbf = availableStats[len(availableStats)/2].End.Sub(availableStats[len(availableStats)/2].Begin)
	} else {
		out.Mtbf = 0
	}
	if len(brokenStats) != 0 {
		out.Mtbr = brokenStats[len(brokenStats)/2].End.Sub(brokenStats[len(brokenStats)/2].Begin)
	} else {
		out.Mtbr = 0
	}
	out.Broken = broken
	out.Total = total
	out.States = stats[0:min(len(stats), 30)]
	if err := json.NewEncoder(w).Encode(&out); err != nil {
		log.Println(err)
	}
}
