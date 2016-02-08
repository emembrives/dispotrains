package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/emembrives/tinkerings/dispotrains.webapp/storage"

	"github.com/eknkc/dateformat"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	homeTmpl         = template.Must(template.ParseFiles("templates/lines.html", "templates/footer.html", "templates/header.html"))
	lineTmpl         = template.Must(template.ParseFiles("templates/line.html", "templates/footer.html", "templates/header.html"))
	stationTmpl      = template.Must(template.ParseFiles("templates/station.html", "templates/footer.html", "templates/header.html"))
	stationStatsTmpl = template.Must(template.ParseFiles("templates/stats.html", "templates/footer.html", "templates/header.html"))
)

type Line struct {
	Network      string
	ID           string
	GoodStations []*storage.Station `bson:"goodStations"`
	BadStations  []*storage.Station `bson:"badStations"`
	LastUpdate   time.Time
}

type LineHolder struct {
	Value Line
}

type LineSlice []LineHolder

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

func (ls LineSlice) Lines() []Line {
	r := make([]Line, len(ls))
	for i, v := range ls {
		r[i] = v.Value
	}
	return r
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-control", "public, max-age=86400")
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("dispotrains").C("lines")
	var lines LineSlice = make(LineSlice, 0)
	c.Find(nil).Sort("value.network", "value.id").All(&lines)
	homeTmpl.Execute(w, lines.Lines())
}

func LineHandler(w http.ResponseWriter, req *http.Request) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("dispotrains").C("lines")

	vars := mux.Vars(req)
	lineId := vars["line"]

	var line LineHolder
	c.Find(bson.M{"_id": lineId}).One(&line)
	w.Header().Set("Last-Modified", line.Value.LastUpdate.UTC().Format(time.RFC1123))
	if err = lineTmpl.Execute(w, line.Value); err != nil {
		log.Fatal(err)
	}
}

func StationHandler(w http.ResponseWriter, req *http.Request) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("dispotrains").C("stations")

	vars := mux.Vars(req)
	stationName := vars["station"]

	var station DisplayStation
	c.Find(bson.M{"name": stationName}).One(&station)
	for _, elevator := range station.Elevators {
		if elevator.Status.State != "Disponible" {
			station.BadElevators++
		}
	}
	w.Header().Set("Last-Modified", station.LastUpdate.UTC().Format(time.RFC1123))
	if err = stationTmpl.Execute(w, station); err != nil {
		log.Fatal(err)
	}
}

func StatsHandler(w http.ResponseWriter, req *http.Request) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	cStations := session.DB("dispotrains").C("stations")
	cStatuses := session.DB("dispotrains").C("statuses")
	cStatistics := session.DB("dispotrains").C("statistics")

	vars := mux.Vars(req)
	stationName := vars["station"]

	var station DisplayStation
	cStations.Find(bson.M{"name": stationName}).One(&station)
	var elevatorIds []string
	for _, stationElevator := range station.Elevators {
		elevatorIds = append(elevatorIds, stationElevator.ID)
	}
	var dbStatuses []dataStatus

	index := mgo.Index{
		Key: []string{"elevator"},
	}
	err = cStatuses.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	cStatuses.Find(bson.M{"elevator": bson.M{"$in": elevatorIds}}).
		Sort("lastupdate").
		All(&dbStatuses)

	events, reports := statusesToEvents(dbStatuses)

	var stats storage.StationStats
	err = cStatistics.Find(bson.M{"name": stationName}).One(&stats)
	if err != nil {
		panic(err)
	}

	type TemplateData struct {
		Station   DisplayStation
		Events    map[string][]string
		Reports   []string
		StartDate string
		EndDate   string
		Stats     storage.StationStats
	}
	var templateData TemplateData
	if len(elevatorIds) != 0 {
		templateData = TemplateData{
			station, events, reports, reports[0], reports[len(reports)-1], stats}
	} else {
		templateData = TemplateData{station, events, reports, "", "", stats}
	}
	if err = stationStatsTmpl.Execute(w, templateData); err != nil {
		log.Fatal(err)
	}
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

// VoronoiHandler sends historical data for the Voronoi map.
func VoronoiHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
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
	if err = json.NewEncoder(w).Encode(&jsonData); err != nil {
		log.Println(err)
	}
}

func AppHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("dispotrains").C("stations")
	var stations []bson.M = make([]bson.M, 0)
	if err := c.Find(nil).All(&stations); err != nil {
		log.Println(err)
	}
	jsonStations := make([]bson.M, 0)
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/ligne/{line}", LineHandler)
	r.HandleFunc("/gare/{station}", StationHandler)
	r.HandleFunc("/gare/{station}/stats", StatsHandler)
	r.HandleFunc("/app/GetStations/", AppHandler)
	r.HandleFunc("/app/AllStats/", VoronoiHandler)
	r.PathPrefix("/static/").Handler(CacheRequest(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
