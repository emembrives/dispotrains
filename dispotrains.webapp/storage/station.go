package storage

import (
	"net/url"
	"strings"
	"time"
)

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type Station struct {
	DisplayName  string
	Name         string
	City         string
	Position     Coordinates
	OsmID        string
	Lines        []*Line
	Elevators    []*Elevator
	code         string
	HasElevators bool
	LastUpdate   time.Time
}

func NewStation(name, city, code string) *Station {
	var station *Station = new(Station)
	station.Name = strings.Replace(strings.TrimSpace(name), "/", "-", -1)
	station.DisplayName = computeDisplayName(station.Name)
	station.City = strings.TrimSpace(city)
	station.code = code
	station.HasElevators = true
	return station
}

func computeDisplayName(name string) string {
	var displayName string = strings.Title(strings.ToLower(name))
	if displayName != "Gare De Lyon" && displayName != "Gare Du Nord" {
		displayName = strings.Replace(displayName, "Gare De ", "", 1)
		displayName = strings.Replace(displayName, "Gare Du ", "", 1)
	}
	displayName = strings.Replace(displayName, " De ", " de ", -1)
	displayName = strings.Replace(displayName, " Du ", " du ", -1)
	displayName = strings.Replace(displayName, " D ", " d'", -1)
	displayName = strings.Replace(displayName, " D'", " d'", -1)
	displayName = strings.Replace(displayName, " Le ", " le ", -1)
	displayName = strings.Replace(displayName, " Les ", " les ", -1)
	displayName = strings.Replace(displayName, " La ", " la ", -1)
	displayName = strings.Replace(displayName, " L ", " l'", -1)
	displayName = strings.Replace(displayName, " En ", " en ", -1)
	displayName = strings.Replace(displayName, " Au ", " au ", -1)
	displayName = strings.Replace(displayName, " Aux ", " aux ", -1)
	displayName = strings.Replace(displayName, " A ", " à ", -1)
	displayName = strings.Replace(displayName, " Sur ", " sur ", -1)
	displayName = strings.Replace(displayName, " Sous ", " sous ", -1)
	return displayName
}

func NewRampStation(name, city string) *Station {
	var station *Station = new(Station)
	station.Name = strings.Replace(strings.TrimSpace(name), "/", "-", -1)
	station.DisplayName = computeDisplayName(station.Name)
	station.City = strings.TrimSpace(city)
	station.HasElevators = false
	return station
}

func (station *Station) AttachLine(line *Line) {
	station.Lines = append(station.Lines, line)
	line.attachStation(station)
}

func (station *Station) NewElevator(id, situation, direction string) *Elevator {
	var elevator *Elevator = new(Elevator)
	elevator.ID = strings.TrimSpace(id)
	elevator.Situation = strings.TrimSpace(situation)
	elevator.Direction = strings.TrimSpace(direction)
	elevator.station = station
	station.Elevators = append(station.Elevators, elevator)
	return elevator
}

func (station *Station) GetElevators() []*Elevator {
	return station.Elevators
}

func (station *Station) GetURL() *url.URL {
	var path *url.URL
	path, err := url.Parse("http://www.infomobi.com/fr/voyageurs-en-fauteuil/transports-publics-accessibles/gares-et-stations-accessibles/disponibilite-des-ascenseurs/?tx_stifinfomobi_pi3[externalcode]=changeme")
	if err != nil {
		panic(err)
	}
	q := path.Query()
	q.Set("tx_stifinfomobi_pi3[externalcode]", station.code)
	path.RawQuery = q.Encode()
	return path
}
