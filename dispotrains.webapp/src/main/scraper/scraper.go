package main

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/client"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/zabawaba99/firego.v1"
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
        if (this["update"] == null || status.lastupdate > this["update"]) {
            this["update"] = status.lastupdate;
        }
        if (status.state == "Disponible") {
            continue;
        } else {
            this["status"] = false;
        }
    }
    delete this["elevators"];

	if (lines.length == 0) {
		throw "No lines for station " + this;
	}

    for (var i = 0; i < lines.length; i++) {
		var key = {"network": lines[i].network, "id": lines[i].id};
        var line = {"network": lines[i].network, "id": lines[i].id};
        line.update = this.update;
        line.goodStations = [];
        line.badStations = [];
        if (this.status) {
            line.goodStations = [this];
        } else {
            line.badStations = [this];
        }
        emit(key, line);
    }
}`

const reduceLines string = `function(key, lines) {
	var line = {"network": key.network, "id": key.id};
	line["update"] = null;
	line.goodStations = [];
	line.badStations = [];
    for (var idx = 0; idx < lines.length; idx++) {
		  var currentLine = lines[idx];
		  for (var i = 0; i < currentLine.badStations.length; i++) {
		  	line.badStations.push(currentLine.badStations[i]);
		  }
		  for (var i = 0; i < currentLine.goodStations.length; i++) {
              line.goodStations.push(currentLine.goodStations[i]);
		  }

		  if (line["update"] == null || line.update < lines[idx].update) {
		  	line.update = lines[idx].update;
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

func uploadToFirebase(session *mgo.Session) error {
	d, err := ioutil.ReadFile("/data/key/dispotrains.json")
	if err != nil {
		return err
	}

	conf, err := google.JWTConfigFromJSON(d, "https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/firebase.database")
	if err != nil {
		return err
	}

	fb := firego.New("https://dispotrains-bbaaa.firebaseio.com", conf.Client(oauth2.NoContext))

	c := session.DB("dispotrains").C("stations")
	var stations []bson.M
	if err := c.Find(nil).All(&stations); err != nil {
		log.Println(err)
	}
	var jsonStations []bson.M
	for _, station := range stations {
		delete(station, "_id")
		jsonStations = append(jsonStations, station)
	}
	return fb.Set(jsonStations)
}

func main() {
	session, err := mgo.DialWithTimeout("db", time.Minute)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior
	session.SetMode(mgo.Monotonic, true)

	// Retrieve lines from STIF website.
	c := session.DB("dispotrains").C("stations")
	lines, err := client.GetAllLines()
	if err != nil {
		panic(err)
	}

	// Build the station database collection.
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
	stationIndex := mgo.Index{
		Key:        []string{"name"},
		Background: true,
		DropDups:   true,
		Unique:     true,
	}
	err = c.EnsureIndex(stationIndex)
	if err != nil {
		panic(err)
	}
	for _, station := range stations {
		_, err = c.Upsert(bson.M{"name": station.Name}, &station)
		if err != nil {
			panic(err)
		}
	}

	// Build the lines database collection.
	job := &mgo.MapReduce{
		Map:    mapStationsToLines,
		Reduce: reduceLines,
		Out:    bson.M{"replace": "lines"},
	}
	_, err = c.Find(bson.M{}).MapReduce(job, nil)
	if err != nil {
		panic(err)
	}

	// Append the new statuses to the database log.
	c = session.DB("dispotrains").C("statuses")
	index := mgo.Index{
		Key:        []string{"state", "lastupdate", "elevator"},
		Background: true,
		DropDups:   true,
		Sparse:     true,
		Unique:     true,
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
	uploadToFirebase(session)
}
