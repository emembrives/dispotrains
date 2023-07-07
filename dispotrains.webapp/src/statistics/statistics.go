package statistics

import (
	"bytes"
	"log"
	"time"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
	"github.com/linxGnu/grocksdb"
	"github.com/vmihailenco/msgpack/v5"
)

type NetworkStats struct {
	Good    int
	Bad     int
	LongBad int
}

func newElevatorState(status storage.StatusStorage) *storage.ElevatorState {
	state := &storage.ElevatorState{}
	state.Elevator = status.ElevatorID
	state.State = status.State
	state.Begin = status.LastUpdate
	state.End = status.LastUpdate
	return state
}

func RecomputeElevatorStatistics(db *grocksdb.DB) error {
	wo := grocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()

	statsWriteBatch := grocksdb.NewWriteBatch()
	defer statsWriteBatch.Destroy()
	statsWriteBatch.DeleteRange([]byte(storage.BucketStatistics), storage.LastKeyOfBucket(storage.BucketStatistics))
	if err := db.Write(wo, statsWriteBatch); err != nil {
		return err
	}
	fo := grocksdb.NewDefaultFlushOptions()
	defer fo.Destroy()
	if err := db.Flush(fo); err != nil {
		log.Panic(err)
	}
	r := grocksdb.Range{
		Start: []byte{0x00},
		Limit: []byte{0xFF},
	}
	db.CompactRange(r)
	return ComputeElevatorStatistics(db)
}

// ComputeElevatorStatistics computes and stores per-elevator statistics.
func ComputeElevatorStatistics(db *grocksdb.DB) error {
	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()

	ro.SetFillCache(false)
	iterator := db.NewIterator(ro)
	defer iterator.Close()

	iterator.Seek([]byte(storage.BucketStatuses))
	for iterator.ValidForPrefix([]byte(storage.BucketStatuses)) {
		statsWriteBatch := grocksdb.NewWriteBatch()
		defer statsWriteBatch.Destroy()
		k := iterator.Key().Data()
		elevatorId := string(bytes.Split(k, []byte("/"))[1])

		var elevatorState *storage.ElevatorState

		elevatorStatisticsIterator := db.NewIterator(ro)
		defer elevatorStatisticsIterator.Close()

		elevatorStatisticsIterator.SeekForPrev(storage.LastKeyOfBucket(storage.BucketStatistics, elevatorId))
		if elevatorStatisticsIterator.ValidForPrefix(storage.MakeKey(storage.BucketStatistics, elevatorId)) {
			elevatorState = &storage.ElevatorState{}
			if err := msgpack.Unmarshal(elevatorStatisticsIterator.Value().Data(), elevatorState); err != nil {
				return err
			}
		}

		var query []byte
		if elevatorState != nil {
			query = storage.LastKeyOfBucket(storage.BucketStatuses, elevatorId, elevatorState.End.Format(time.RFC3339))
			iterator.Seek(query)
		}

		var status storage.StatusStorage
		for ; iterator.ValidForPrefix(storage.MakeKey(storage.BucketStatuses, elevatorId)); iterator.Next() {
			if err := msgpack.Unmarshal(iterator.Value().Data(), &status); err != nil {
				return err
			}
			if elevatorState == nil {
				elevatorState = newElevatorState(status)
				continue
			}
			if status.LastUpdate.Before(elevatorState.Begin) {
				continue
			}
			elevatorState.End = status.LastUpdate
			if _, isUnknown := storage.UnknownStates[status.State]; isUnknown {
				continue
			}
			if (status.State == "Disponible") != (elevatorState.State == "Disponible") {
				b, err := msgpack.Marshal(elevatorState)
				if err != nil {
					return err
				}
				statsWriteBatch.Put(storage.MakeKey(storage.BucketStatistics, elevatorId, elevatorState.Begin.Format(time.RFC3339)), b)
				elevatorState = newElevatorState(status)
			}
		}
		b, err := msgpack.Marshal(elevatorState)
		if err != nil {
			return err
		}
		statsWriteBatch.Put(storage.MakeKey(storage.BucketStatistics, elevatorId, elevatorState.Begin.Format(time.RFC3339)), b)

		wo := grocksdb.NewDefaultWriteOptions()
		defer wo.Destroy()
		if err := db.Write(wo, statsWriteBatch); err != nil {
			return err
		}
		// Go to the next elevator.
	}
	return nil
}

func ComputeGlobalStatistics(db *grocksdb.DB) (*NetworkStats, error) {
	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()

	iterator := db.NewIterator(ro)
	defer iterator.Close()

	ns := NetworkStats{}
	longLimit := time.Now().AddDate(0, 0, -3)
	limit := time.Now().AddDate(0, 0, -10)

	for iterator.Seek([]byte(storage.BucketStatistics)); iterator.ValidForPrefix([]byte(storage.BucketStatistics)); iterator.Next() {
		k := iterator.Key().Data()
		elevatorId := string(bytes.Split(k, []byte("/"))[1])

		lastStateQuery := storage.LastKeyOfBucket(storage.BucketStatistics, elevatorId)
		iterator.SeekForPrev(lastStateQuery)

		var state storage.ElevatorState
		if err := msgpack.Unmarshal(iterator.Value().Data(), &state); err != nil {
			return nil, err
		}

		if state.End.Before(limit) {
			continue
		} else if state.State == "Disponible" {
			ns.Good++
		} else {
			ns.Bad++
			if state.Begin.Before(longLimit) {
				ns.LongBad++
			}
		}
	}
	return &ns, nil
}
