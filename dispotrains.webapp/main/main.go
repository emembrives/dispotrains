package main

import (
	"strings"

	"github.com/emembrives/tinkerings/dispotrains.webapp/client"
	"github.com/emembrives/tinkerings/dispotrains.webapp/storage"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const mapStationsToLines string = `function() {
    var lines = this.lines;
    delete this["lines"];
    delete this["_id"];
    this["status"] = true;
    this["update"] = null;

    for (var i = 0; i < this.elevators.length; i++) {
        var elevator = this.elevators[i];
        var status = elevator.status;
        if (this["update"] != null && status.lastupdate > this["update"]) {
            this["update"] = status.lastupdate;
        } else {
            this["update"] = status.lastupdate;
        }
        if (status.state == "Disponible") {
            continue;
        } else {
            this["status"] = false;
            break;
        }
    }
    delete this["elevators"];

    for (var i = 0; i < lines.length; i++) {
        var line = lines[i];
        delete line["_id"];
        line.update = this.update;
        line.goodStations = [];
        line.badStations = [];
        if (this.status) {
            line.goodStations = [this];
        } else {
            line.badStations = [this];
        }
        emit(line.id, line);
    }
}`

const reduceLines string = `function(keySKU, lines) {
    var line = lines[0];
    for (var idx = 1; idx < lines.length; idx++) {
        if (lines[idx].badStations.length > 0) {
            line.badStations.push(lines[idx].badStations[0]);
            if (line.update < lines[idx].badStations[0].update) {
                line.update = lines[idx].badStations[0].update;
            }
        }
        if (lines[idx].goodStations.length > 0) {
            line.goodStations.push(lines[idx].goodStations[0]);
            if (line.update < lines[idx].goodStations[0].update) {
                line.update = lines[idx].goodStations[0].update;
            }
        }
    }
    var sortFunc = function(a, b) {
        if (a.displayname < b.displayname) {
            return -1;
        } else if (a.displayname == b.displayname) {
            return 0;
        } else {
            return 1;
        }};
    line.badStations.sort(sortFunc);
    line.goodStations.sort(sortFunc);
    return line;
}`

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior
	// session.SetMode(mgo.Monotonic, true)

	c := session.DB("dispotrains").C("stations")
	lines, err := client.GetAllLines()

	if err != nil {
		panic(err)
	}
	stations := make(map[string]*storage.Station)
	for _, line := range lines {
		for _, station := range line.GetStations() {
			if _, ok := stations[strings.ToLower(station.Name)]; ok == true {
				for _, sLine := range station.Lines {
					stations[strings.ToLower(station.Name)].AttachLine(sLine)
				}
			} else {
				stations[strings.ToLower(station.Name)] = station
			}
		}
	}
	AddPositionToStations(stations)

	c.RemoveAll(bson.M{})
	for _, station := range stations {
		err = c.Insert(&station)
		if err != nil {
			panic(err)
		}
	}
	job := &mgo.MapReduce{Map: mapStationsToLines, Reduce: reduceLines, Out: bson.M{"replace": "lines"}}
	_, err = c.Find(nil).MapReduce(job, nil)
	if err != nil {
		panic(err)
	}

	c = session.DB("dispotrains").C("statuses")
	index := mgo.Index{
		Key:        []string{"state", "lastupdate", "elevator"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	for _, station := range stations {
		for _, elevator := range station.GetElevators() {
			bsonStatus := bson.M{
				"state":      elevator.Status.State,
				"lastupdate": elevator.Status.LastUpdate,
				"elevator":   elevator.ID,
			}
			err = c.Insert(bsonStatus)
			if err != nil && !mgo.IsDup(err) {
				panic(err)
			}
		}
	}
}

/*
		for _, station := range line.GetStations() {
			if !station.HasElevators {
				continue
			}
			for _, elevator := range station.GetElevators() {
				status := elevator.GetLastStatus()
			}
		}
	}
}*/
