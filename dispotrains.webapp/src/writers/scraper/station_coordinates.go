package scraper

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
)

type csvStation struct {
	Name, City          string
	Latitude, Longitude float64
	OsmID               string
}

func newCsvStation(name, city, latitude, longitude, osmid string) csvStation {
	convLatitude, err := strconv.ParseFloat(latitude, 64)
	panicOnError(err)
	convLongitude, err := strconv.ParseFloat(longitude, 64)
	panicOnError(err)
	return csvStation{name, city, convLatitude, convLongitude, osmid}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func readCsvFile() (stations []csvStation) {
	f, err := os.Open("stations-coordinates.csv")
	panicOnError(err)

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	panicOnError(err)

	for _, record := range records[1:] {
		stations = append(stations, newCsvStation(record[0], record[1], record[2], record[3], record[4]))
	}
	return stations
}

// addPositionToStations adds Latitude/Longitude coordinates data from side CSV
// file to the station map.
func addPositionToStations(stations []*storage.Station) {
	csvStations := readCsvFile()

	stationMap := make(map[string]*storage.Station)
	for _, station := range stations {
		stationMap[strings.ToLower(station.Name)] = station
	}

	for _, csvStation := range csvStations {
		if _, ok := stationMap[strings.ToLower(csvStation.Name)]; !ok {
			continue
		}

		stationMap[strings.ToLower(csvStation.Name)].Position.Latitude = csvStation.Latitude
		stationMap[strings.ToLower(csvStation.Name)].Position.Longitude = csvStation.Longitude
		stationMap[strings.ToLower(csvStation.Name)].OsmID = csvStation.OsmID
	}
}
