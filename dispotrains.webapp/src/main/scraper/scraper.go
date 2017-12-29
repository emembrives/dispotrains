package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/assistant"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/client"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/environment"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/zabawaba99/firego.v1"
)

func uploadToFirebase(session *mgo.Session) error {
	d, err := ioutil.ReadFile("/dispotrains/key/dispotrains.json")
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
	session, err := mgo.DialWithTimeout(environment.GetMongoDbAddress(), time.Minute)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Retrieve lines from STIF website.
	c := session.DB("dispotrains").C("stations")
	lines, stations, err := client.GetAndParseLines()
	if err != nil {
		panic(err)
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
	// Build the lines database collection.
	bulk := c.Bulk()
	bulk.RemoveAll(nil)
	for _, station := range stations {
		bulk.Insert(station)
	}
	_, err = bulk.Run()
	if err != nil {
		panic(err)
	}

	// Build the lines database collection.
	c = session.DB("dispotrains").C("lines")
	bulk = c.Bulk()
	bulk.RemoveAll(nil)
	for _, line := range lines {
		bulk.Insert(bson.M{
			"network":      line.Network,
			"id":           line.ID,
			"lastupdate":   line.LastUpdate,
			"goodstations": line.GoodStations(),
			"badstations":  line.BadStations(),
		})
	}
	_, err = bulk.Run()
	if err != nil {
		panic(err)
	}

	// Append the new statuses to the database log.
	c = session.DB("dispotrains").C("statuses")
	index := mgo.Index{
		Key:        []string{"state", "lastupdate", "elevator"},
		Background: true,
		Sparse:     true,
		Unique:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	bulk = c.Bulk()
	for _, station := range stations {
		for _, elevator := range station.GetElevators() {
			if elevator.Status == nil {
				continue
			}
			bsonStatus := bson.M{
				"state":      elevator.Status.State,
				"lastupdate": elevator.Status.LastUpdate,
				"elevator":   elevator.ID,
			}
			if elevator.Status.Forecast != nil {
				bsonStatus["forecast"] = elevator.Status.Forecast
			}
			bulk.Insert(bsonStatus)
		}
	}
	bulk.Unordered()
	_, err = bulk.Run()
	if err != nil && !mgo.IsDup(err) {
		panic(err)
	}
	assistant.UpdateStationList(stations)
}
