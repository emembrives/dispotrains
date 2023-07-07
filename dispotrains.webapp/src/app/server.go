package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
	"github.com/linxGnu/grocksdb"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/gorilla/mux"
)

var (
	session *grocksdb.DB = nil
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

func createSessionOrDie() *grocksdb.DB {
	db, err := storage.GetDatabase()
	if err != nil {
		log.Panicf("Unable to get database: %v", err)
	}
	return db
}

func GetStationsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	iter := session.NewIterator(ro)
	defer iter.Close()

	var stations []*storage.Station

	for iter.Seek(storage.MakeKey(storage.BucketStations)); iter.ValidForPrefix(storage.MakeKey(storage.BucketStations)); iter.Next() {
		var station storage.Station
		err := msgpack.Unmarshal(iter.Value().Data(), &station)
		if err != nil {
			log.Panic(err)
		}
		stations = append(stations, &station)
	}
	json.NewEncoder(w).Encode(&stations)
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
	flag.Parse()

	session = createSessionOrDie()
	defer session.Close()

	go control()

	r := mux.NewRouter()
	r.HandleFunc("/app/GetStations/", GetStationsHandler)
	r.HandleFunc("/app/netStats/", NetworkStatsHandler)
	r.HandleFunc("/app/Elevator/{id}", ElevatorHandler)
	r.PathPrefix("/static/").Handler(CacheRequest(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))
	r.PathPrefix("/").Handler(CacheRequest(http.FileServer(NewFileHandlerWithDefault("index.html", "web.v2"))))
	http.Handle("/", r)
	log.Panic(http.ListenAndServe("0.0.0.0:9000", nil))
}
