package storage

import (
	"bytes"
	"flag"
	"log"
	"strings"

	"github.com/linxGnu/grocksdb"
)

const (
	BucketStations   = "stations"
	BucketStatuses   = "statuses"
	BucketStatistics = "statistics"
)

var (
	databasePath = flag.String("database_path", "", "Path to the database")
)

/*
 * Database scheme:
 * stations/STATION_NAME/: storage.Station
 * statuses/ELEVATOR_ID/LAST_UPDATE_ISO/: storage.StorageStatus
 * statistics/ELEVATOR_ID/RANGE_START_ISO/: storage.ElevatorState
 *
 */

func GetDatabase() (*grocksdb.DB, error) {
	if !flag.Parsed() || *databasePath == "" {
		log.Panic("--database_path must be provided")
	}

	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	// bbto.SetBlockCache(grocksdb.NewLRUCache(3 << 30))
	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := grocksdb.OpenDb(opts, *databasePath)
	return db, err
}

func MakeKey(parts ...string) []byte {
	buf := bytes.NewBuffer(nil)
	if _, err := buf.WriteString(strings.Join(parts, "/")); err != nil {
		log.Panic(err)
	}
	buf.WriteString("/")
	return buf.Bytes()
}

func LastKeyOfBucket(parts ...string) []byte {
	return append(MakeKey(parts...), byte(0xFF))
}
