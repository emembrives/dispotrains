package main

import (
	"crypto/elliptic"
	"encoding/json"
	"log"
	"net/http"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/push"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
)

func PushRegisterHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if req.Method != "POST" {
		return
	}

	var data storage.Subscription
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	push.Register(session, data)
}

func PushToAllHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if req.Method != "GET" {
		return
	}
	push.PushToAll(session)
}

func GetVAPIDKeyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-control", "public, max-age=86400")

	key := push.GetOrCreateVAPIDKey(session).ToElliptic()
	json.NewEncoder(w).Encode(elliptic.Marshal(key.Curve, key.X, key.Y))
}
