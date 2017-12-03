package main

import (
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/assistant"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
)

func main() {
	station := storage.Station{
		Name:        "GARE DE LA DEFENSE",
		DisplayName: "Gare de La Défense",
	}
	stations := make([]*storage.Station, 1)
	stations[0] = &station
	assistant.UpdateStationList(stations)
}
