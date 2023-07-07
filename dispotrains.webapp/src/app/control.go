package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/statistics"
	jsondump "github.com/emembrives/dispotrains/dispotrains.webapp/src/writers/json_dump"
	loaddump "github.com/emembrives/dispotrains/dispotrains.webapp/src/writers/load_dump"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/writers/scraper"
	"github.com/gorilla/mux"
)

var (
	dumpPath = flag.String("dump_path", "", "Path to BZip2 dump file")
)

func control() {
	r := mux.NewRouter()
	r.HandleFunc("/load", func(w http.ResponseWriter, r *http.Request) {
		loaddump.Loaddump(session, *dumpPath)
	})
	r.HandleFunc("/dump", func(w http.ResponseWriter, r *http.Request) {
		jsondump.Jsondump(session, *dumpPath)
	})
	r.HandleFunc("/scrape", func(w http.ResponseWriter, r *http.Request) {
		scraper.Scraper(session)
	})
	r.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		if err := statistics.RecomputeElevatorStatistics(session); err != nil {
			log.Panic(err)
		}
	})
	serveMux := http.NewServeMux()
	serveMux.Handle("/", r)
	log.Panic(http.ListenAndServe("127.0.0.1:9001", serveMux))
}
