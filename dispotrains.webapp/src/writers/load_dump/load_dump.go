package loaddump

import (
	"archive/tar"
	"bufio"
	"compress/bzip2"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
	"github.com/linxGnu/grocksdb"
	"github.com/vmihailenco/msgpack/v5"
)

type mgoStatus struct {
	State      string `json:"state"`
	Elevator   string `json:"elevator"`
	LastUpdate struct {
		Date time.Time `json:"$date"`
	} `json:"lastupdate"`
	Forecast *struct {
		Date time.Time `json:"$date"`
	} `json:"forecast"`
}

func Loaddump(db *grocksdb.DB, dumpPath string) {
	flag.Parse()

	file, err := os.Open(path.Join(dumpPath, "dump.tar.bz2"))
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	bzip2Reader := bzip2.NewReader(file)
	tarReader := tar.NewReader(bzip2Reader)

	for {
		header, err := tarReader.Next()
		if err != nil {
			log.Panic(err)
		}
		if header.Name == "statuses.json" {
			break
		}
	}

	scanner := bufio.NewScanner(tarReader)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	wo := grocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()
	wo.DisableWAL(true)
	for scanner.Scan() {
		var status mgoStatus
		err := json.Unmarshal(scanner.Bytes(), &status)
		if err != nil {
			log.Panicf("Error while parsing %s: %v", scanner.Text(), err)
		}
		storageStatus := storage.StatusStorage{
			State:      status.State,
			LastUpdate: status.LastUpdate.Date,
			ElevatorID: status.Elevator,
		}
		if status.Forecast != nil {
			storageStatus.Forecast = &status.Forecast.Date
		} else if strings.Contains(storageStatus.State, " jusqu'au ") {
			forecastStr := storageStatus.State[len(storageStatus.State)-10:]
			t, err := time.Parse("02/01/2006", forecastStr)
			if err != nil {
				log.Panicf("Unable to parse date in %s", storageStatus.State)
				continue
			}
			storageStatus.Forecast = &t
		}
		bytes, err := msgpack.Marshal(storageStatus)
		if err != nil {
			log.Panic(err)
		}
		db.Put(wo, storage.MakeKey(storage.BucketStatuses, storageStatus.ElevatorID, storageStatus.LastUpdate.Format(time.RFC3339)), bytes)
	}

	r := grocksdb.Range{
		Start: []byte{0x00},
		Limit: []byte{0xFF},
	}
	db.CompactRange(r)
}
