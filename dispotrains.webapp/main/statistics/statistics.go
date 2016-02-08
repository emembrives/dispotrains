package main

import (
	"sort"
	"time"

	"github.com/emembrives/tinkerings/dispotrains.webapp/storage"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type dataStatus struct {
	Elevator   string
	State      string
	Lastupdate time.Time
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cStations := session.DB("dispotrains").C("stations")
	cStatuses := session.DB("dispotrains").C("statuses")
	cStatistics := session.DB("dispotrains").C("statistics")

	err = cStatuses.EnsureIndexKey("elevator")
	if err != nil {
		panic(err)
	}
	err = cStatistics.EnsureIndexKey("name")
	if err != nil {
		panic(err)
	}

	stations := make([]storage.Station, 0, 0)
	err = cStations.Find(bson.M{}).All(&stations)
	if err != nil {
		panic(err)
	}

	for _, station := range stations {
		func(currentStation storage.Station) {
			stationStats := statsHandler(currentStation, cStatuses)
			_, err := cStatistics.Upsert(bson.M{"name": currentStation.Name}, stationStats)
			if err != nil {
				panic(err)
			}
		}(station)
	}
}

func statsHandler(station storage.Station, cStatuses *mgo.Collection) storage.StationStats {
	var elevatorIds []string
	for _, stationElevator := range station.Elevators {
		elevatorIds = append(elevatorIds, stationElevator.ID)
	}
	var dbStatuses []dataStatus

	cStatuses.Find(bson.M{"elevator": bson.M{"$in": elevatorIds}}).
		Sort("lastupdate").
		All(&dbStatuses)

	events, reports := statusesToEvents(dbStatuses)
	stats := statusesToStatistics(events, reports, dbStatuses)
	stats.Name = station.Name
	return stats
}

func statusesToEvents(dbStatuses []dataStatus) (map[string][]string, []string) {
	events := make(map[string][]string)
	reportSet := make(map[string]bool)
	for _, status := range dbStatuses {
		dateStr := status.Lastupdate.Format(time.RFC3339)
		reportSet[dateStr] = true
		if _, ok := events[status.Elevator]; !ok {
			events[status.Elevator] = make([]string, 0)
		}
		if status.State != "Disponible" {
			events[status.Elevator] = append(events[status.Elevator], dateStr)
		}
	}

	reports := make(sort.StringSlice, 0, len(reportSet))
	for key := range reportSet {
		reports = append(reports, key)
	}
	reports.Sort()
	return events, reports
}

func statusesToStatistics(events map[string][]string, reports []string, dbStatuses []dataStatus) storage.StationStats {
	stats := storage.StationStats{}
	stats.Reports = len(reports)
	reportDays := make(map[string]bool)
	for _, date := range reports {
		reportDays[date[0:10]] = true
	}
	stats.ReportDays = len(reportDays)
	stats.Elevators = make(map[string]storage.ElevatorStats)
	malfunctionDays := make(map[string]bool)
	for elevatorName, statusDates := range events {
		elevatorStats := storage.ElevatorStats{
			Name:         elevatorName,
			Malfunctions: len(statusDates),
		}
		stats.Malfunctions += len(statusDates)
		malfunctionElevatorDays := make(map[string]bool)
		for _, date := range statusDates {
			malfunctionDays[date[0:10]] = true
			malfunctionElevatorDays[date[0:10]] = true
		}
		elevatorStats.MalfunctionDays = len(malfunctionElevatorDays)
		elevatorStats.FunctionDays = len(reportDays) - len(malfunctionElevatorDays)
		if len(reportDays) != 0 {
			elevatorStats.PercentFunction = float64(elevatorStats.FunctionDays) * 100 / float64(len(reportDays))
		} else {
			elevatorStats.PercentFunction = 100.0
		}
		stats.Elevators[elevatorName] = elevatorStats
	}
	stats.MalfunctionDays = len(malfunctionDays)
	stats.FunctionDays = len(reportDays) - len(malfunctionDays)
	if len(reportDays) != 0 {
		stats.PercentFunction = float64(stats.FunctionDays) * 100 / float64(len(reportDays))
	} else {
		stats.PercentFunction = 100.0
	}
	return stats
}
