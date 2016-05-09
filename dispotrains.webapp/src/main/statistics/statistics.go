package main

import (
	"sort"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
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

// Events compile events for a station.
type Events struct {
	eventsByElevator map[string][]string
	reportDays       []string
	firstDay         time.Time
	lastDay          time.Time
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

	events := statusesToEvents(dbStatuses)
	stats := statusesToStatistics(events)
	stats.Name = station.Name
	return stats
}

func statusesToEvents(dbStatuses []dataStatus) Events {
	events := make(map[string][]string)
	reportSet := make(map[string]bool)
	var firstStatus, lastStatus time.Time
	if len(dbStatuses) != 0 {
		firstStatus, lastStatus = dbStatuses[0].Lastupdate, dbStatuses[0].Lastupdate
	}
	for _, status := range dbStatuses {
		if status.State == "Information non disponible" {
			continue
		}
		if lastStatus.Before(status.Lastupdate) {
			lastStatus = status.Lastupdate
		}
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
	return Events{events, reports, firstStatus, lastStatus}
}

func statusesToStatistics(events Events) storage.StationStats {
	stats := storage.StationStats{}
	stats.Reports = len(events.reportDays)
	// This is probably not correct due to daylight saving changes.
	stats.DisplayDays = int(events.lastDay.Truncate(24*time.Hour).Sub(
		events.firstDay.Truncate(24*time.Hour)).Hours() / 24)
	reportDays := make(map[string]bool)
	for _, date := range events.reportDays {
		reportDays[date[0:10]] = true
	}
	stats.ReportDays = len(reportDays)
	stats.Elevators = make(map[string]storage.ElevatorStats)
	malfunctionDays := make(map[string]bool)
	for elevatorName, statusDates := range events.eventsByElevator {
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
