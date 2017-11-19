package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"

	"github.com/eknkc/dateformat"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoDbHost = "db"
)

var (
	session = createSessionOrDie()
)

type Line struct {
	Network      string
	ID           string
	GoodStations []*storage.Station
	BadStations  []*storage.Station
	LastUpdate   time.Time
}

type LineSlice []Line

type DisplayStation struct {
	Name         string
	DisplayName  string
	City         string
	Position     storage.Coordinates
	OsmID        string
	Elevators    []*LocElevator
	LastUpdate   time.Time
	BadElevators int
}

type LocElevator storage.Elevator

type dataStatus struct {
	Elevator   string
	State      string
	Lastupdate time.Time
}

func (e *LocElevator) LocalStatusDate() string {
	return dateformat.FormatLocale(e.Status.LastUpdate, "ddd D MMM à HH:MM", dateformat.French)
}

func createSessionOrDie() *mgo.Session {
	session, err := mgo.Dial(mongoDbHost)
	if err != nil {
		panic(err)
	}
	return session
}

// VoronoiHandler sends historical data for the Voronoi map.
func VoronoiHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	cStatistics := session.DB("dispotrains").C("statistics")

	var stats []bson.M = make([]bson.M, 0)
	if err := cStatistics.Find(nil).All(&stats); err != nil {
		log.Println(err)
	}
	var jsonData []bson.M = make([]bson.M, 0)
	for _, stat := range stats {
		delete(stat, "_id")
		jsonData = append(jsonData, stat)
	}
	if err := json.NewEncoder(w).Encode(&jsonData); err != nil {
		log.Println(err)
	}
}

func GetLinesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-control", "public, max-age=86400")

	c := session.DB("dispotrains").C("lines")
	var lines = make(LineSlice, 0)
	if err := c.Find(nil).Sort("network", "id").All(&lines); err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(&lines)
}

func GetStationsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

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
	json.NewEncoder(w).Encode(&jsonStations)
}

func CacheRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-control", "public, max-age=259200")
		h.ServeHTTP(w, r)
	})
}

type FileHandlerWithDefault struct {
	filename string
	dir      http.Dir
}

func NewFileHandlerWithDefault(filename, path string) *FileHandlerWithDefault {
	return &FileHandlerWithDefault{filename, http.Dir(path)}
}

func (s *FileHandlerWithDefault) Open(name string) (http.File, error) {
	if f, err := s.dir.Open(name); err == nil {
		return f, err
	}
	return s.dir.Open(s.filename)
}

func main() {
	defer session.Close()
	r := mux.NewRouter()
	r.HandleFunc("/app/GetLines/", GetLinesHandler)
	r.HandleFunc("/app/GetStations/", GetStationsHandler)
	r.HandleFunc("/app/AllStats/", VoronoiHandler)
	r.PathPrefix("/static/").Handler(CacheRequest(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))
	r.PathPrefix("/").Handler(CacheRequest(http.FileServer(NewFileHandlerWithDefault("index.html", "dist"))))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("0.0.0.0:9000", nil))
}
