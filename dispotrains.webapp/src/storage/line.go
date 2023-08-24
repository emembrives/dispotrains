package storage

import (
	"net/url"
	"time"
)

type Line struct {
	Network    string `json:"network"`
	ID         string `json:"id"`
	code       string
	stations   []*Station
	LastUpdate time.Time `json:"lastupdate"`
}

func (line *Line) attachStation(station *Station) {
	line.stations = append(line.stations, station)
}

func NewLine(network, id, code string) *Line {
	var line *Line = new(Line)
	line.Network = network
	line.ID = id
	line.code = code
	line.stations = make([]*Station, 0)
	return line
}

func (line *Line) GetStations() []*Station {
	return line.stations
}

func (line *Line) GoodStations() []*Station {
	good := make([]*Station, 0)
	for _, station := range line.stations {
		if station.Available() {
			good = append(good, station)
		}
	}
	return good
}

func (line *Line) BadStations() []*Station {
	bad := make([]*Station, 0)
	for _, station := range line.stations {
		if !station.Available() {
			bad = append(bad, station)
		}
	}
	return bad
}

func (line *Line) GetURL() *url.URL {
	var path *url.URL
	path, err := url.Parse("http://www.infomobi.com/fr/voyageurs-en-fauteuil/transports-publics-accessibles/gares-et-stations-accessibles/?tx_stifinfomobi_pi7[externalcode]=changeme")
	if err != nil {
		panic(err)
	}
	q := path.Query()
	q.Set("tx_stifinfomobi_pi7[externalcode]", line.code)
	path.RawQuery = q.Encode()
	return path
}
