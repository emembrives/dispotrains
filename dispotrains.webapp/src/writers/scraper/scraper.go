package scraper

import (
	"log"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/client"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/statistics"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
	"github.com/linxGnu/grocksdb"
	"github.com/vmihailenco/msgpack/v5"
)

func Scraper(db *grocksdb.DB) {
	_, stations, err := client.GetAndParseLines()
	if err != nil {
		log.Panic(err)
	}
	addPositionToStations(stations)

	stationWriteBatch := grocksdb.NewWriteBatch()
	defer stationWriteBatch.Destroy()
	stationWriteBatch.DeleteRange([]byte(storage.BucketStations), storage.LastKeyOfBucket(storage.BucketStations))

	for _, station := range stations {
		data, err := msgpack.Marshal(station)
		if err != nil {
			log.Panic(err)
		}
		stationWriteBatch.Put(storage.MakeKey(storage.BucketStations, station.Name), data)
	}
	wo := grocksdb.NewDefaultWriteOptions()
	if err := db.Write(wo, stationWriteBatch); err != nil {
		log.Panic(err)
	}

	statusWriteBatch := grocksdb.NewWriteBatch()
	defer statusWriteBatch.Destroy()
	// Append the new statuses to the database log.
	for _, station := range stations {
		for _, elevator := range station.GetElevators() {
			if elevator.Status == nil {
				continue
			}
			bytes, err := msgpack.Marshal(elevator.Status.ToStorage())
			if err != nil {
				log.Panic(err)
			}
			statusWriteBatch.Put(storage.MakeKey(storage.BucketStatuses, elevator.ID, elevator.Status.LastUpdate.Format(time.RFC3339)), bytes)
		}
	}
	if err := db.Write(wo, statusWriteBatch); err != nil {
		log.Panic(err)
	}

	err = statistics.ComputeElevatorStatistics(db)
	if err != nil {
		log.Panic(err)
	}
}
