package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type VAPIDKey struct {
	PrivateKeyX *big.Int
	PrivateKeyY *big.Int
	PrivateKeyD *big.Int
	PublicKey   []byte
}

func PushSubHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	var data map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Push, endpoint: %v+", data)
}

func PushToAllHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	vapid := getOrCreateVAPIDKey()
	signedJWT, err := createSignedJWT("https://fcm.googleapis.com/fcm/send/not-the-right-url")
	if err != nil {
		log.Printf("Error while creating signed JWT: %v", err)
		return
	}
	encodedKey := base64.RawURLEncoding.EncodeToString(vapid.PublicKey)
	log.Println(fmt.Sprintf("vapid t=%s,k=%s", signedJWT, encodedKey))
}

func createSignedJWT(endpoint string) (string, error) {
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	endpointUrl.Path = ""
	endpointUrl.RawQuery = ""
	endpointUrl.Fragment = ""
	key := getOrCreateVAPIDKey()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{
		Audience:  endpointUrl.String(),
		Subject:   "mailto:foo@bar.fr",
		ExpiresAt: time.Now().Add(time.Duration(12) * time.Hour).Unix(),
	})
	log.Println(key)
	privKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     key.PrivateKeyX,
			Y:     key.PrivateKeyY,
		},
		D: key.PrivateKeyD,
	}
	return token.SignedString(privKey)
}

func VAPIDHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-control", "public, max-age=86400")

	c := session.DB("dispotrains").C("pushKey")
	var keyPair VAPIDKey
	err := c.Find(nil).One(&keyPair)
	log.Println(keyPair)
	if err != nil || keyPair.PrivateKeyX == nil {
		log.Printf("Creating new key pair: %+v", err)
		c.DropCollection()
		keyPair = createKeyPair()
		c.Insert(keyPair)
	}
	json.NewEncoder(w).Encode(&keyPair)
}

func getOrCreateVAPIDKey() VAPIDKey {
	c := session.DB("dispotrains").C("pushKey")
	var keyPair VAPIDKey
	err := c.Find(nil).One(&keyPair)
	if err != nil {
		c.DropCollection()
		keyPair = createKeyPair()
		c.Insert(keyPair)
	}
	return keyPair
}

func createKeyPair() VAPIDKey {
	curve := elliptic.P256()
	// priv
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}
	// pub
	return VAPIDKey{
		PrivateKeyX: priv.X,
		PrivateKeyY: priv.Y,
		PrivateKeyD: priv.D,
		PublicKey:   elliptic.Marshal(priv, priv.X, priv.Y),
	}
}
