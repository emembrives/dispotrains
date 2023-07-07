package jsondump

import (
	"archive/tar"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
	"github.com/linxGnu/grocksdb"
	"github.com/ulikunitz/xz"
	"github.com/vmihailenco/msgpack/v5"
)

func Jsondump(db *grocksdb.DB, dumpPath string) {
	flag.Parse()

	f, err := os.Create(path.Join(dumpPath, "dump.tar.xz"))
	if err != nil {
		log.Panic(f)
	}
	defer f.Close()

	xzWriter, err := xz.NewWriter(f)
	if err != nil {
		log.Panic(err)
	}
	defer xzWriter.Close()
	tarWriter := tar.NewWriter(xzWriter)

	snapshot := db.NewSnapshot()
	defer snapshot.Destroy()

	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	ro.SetSnapshot(snapshot)
	iter := db.NewIterator(ro)

	stationsCounter := &CountingWriter{}
	extractFromDB[storage.Station](iter, storage.MakeKey(storage.BucketStations), stationsCounter)

	stations := tar.Header{
		Name:    "stations.json",
		Size:    int64(stationsCounter.Count),
		ModTime: time.Now(),
		Mode:    0640,
	}
	tarWriter.WriteHeader(&stations)
	extractFromDB[storage.Station](iter, storage.MakeKey(storage.BucketStations), tarWriter)

	statusCounter := &CountingWriter{}
	extractFromDB[storage.StatusStorage](iter, storage.MakeKey(storage.BucketStatuses), statusCounter)

	statuses := tar.Header{
		Name:    "statuses.json",
		Size:    int64(statusCounter.Count),
		ModTime: time.Now(),
		Mode:    0640,
	}
	tarWriter.WriteHeader(&statuses)
	extractFromDB[storage.StatusStorage](iter, storage.MakeKey(storage.BucketStatuses), tarWriter)

	err = tarWriter.Close()
	if err != nil {
		log.Panic(err)
	}
}

type CountingWriter struct {
	Count int
}

func (cw *CountingWriter) Write(b []byte) (int, error) {
	l := len(b)
	cw.Count += l
	return l, nil
}

func extractFromDB[t any](iter *grocksdb.Iterator, keyPrefix []byte, writer io.Writer) {
	for iter.Seek(keyPrefix); iter.ValidForPrefix(keyPrefix); iter.Next() {
		var value t
		msgpack.Unmarshal(iter.Value().Data(), &value)
		b, err := json.Marshal(value)
		if err != nil {
			log.Panic(err)
		}
		writer.Write(b)
		writer.Write([]byte("\n"))
	}
}
