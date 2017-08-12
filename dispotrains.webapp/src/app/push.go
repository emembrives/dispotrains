package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/dgrijalva/jwt-go"
)

type VAPIDKey struct {
	PrivateKeyX *big.Int
	PrivateKeyY *big.Int
	PrivateKeyD *big.Int
}

func (k VAPIDKey) GetBSON() (interface{}, error) {
	return bson.D{{"x", k.PrivateKeyX.Bytes()}, {"y", k.PrivateKeyY.Bytes()}, {"d", k.PrivateKeyD.Bytes()}}, nil
}

func (k *VAPIDKey) SetBSON(raw bson.Raw) error {
	var out bson.M
	if err := raw.Unmarshal(&out); err != nil {
		return err
	}
	var tmp big.Int
	var tmpBytes []byte
	var ok bool
	if tmpBytes, ok = out["x"].([]byte); !ok {
		return errors.New("Unable to convert x")
	}
	k.PrivateKeyX = tmp.SetBytes(tmpBytes)
	if tmpBytes, ok = out["y"].([]byte); !ok {
		return errors.New("Unable to convert y")
	}
	k.PrivateKeyY = tmp.SetBytes(tmpBytes)
	if tmpBytes, ok = out["d"].([]byte); !ok {
		return errors.New("Unable to convert d")
	}
	k.PrivateKeyD = tmp.SetBytes(tmpBytes)
	return nil
}

func (k *VAPIDKey) ToElliptic() *ecdsa.PrivateKey {
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     k.PrivateKeyX,
			Y:     k.PrivateKeyY,
		},
		D: k.PrivateKeyD,
	}
}

type Registration struct {
	Subscription map[string]interface{}
	LastUpdated  time.Time
}

func PushSubHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if req.Method != "POST" {
		return
	}

	var data map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}
	registration := Registration{
		Subscription: data,
		LastUpdated:  time.Now(),
	}

	c := session.DB("dispotrains").C("pushSubscribers")
	_, err = c.Upsert(map[string]interface{}{"subscription": data}, registration)
	if err != nil {
		log.Println(err)
	}
}

func PushToAllHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	vapid := getOrCreateVAPIDKey()
	signedJWT, err := createSignedJWT("https://fcm.googleapis.com/fcm/send/not-the-right-url")
	if err != nil {
		log.Printf("Error while creating signed JWT: %v\n", err)
		return
	}
	private := vapid.ToElliptic()
	encodedKey, err := x509.MarshalPKIXPublicKey(&private.PublicKey)
	if err != nil {
		log.Printf("Error encoding key: %v\n", err)
		return
	}
	b64Key := base64.RawURLEncoding.EncodeToString(encodedKey)
	log.Println(fmt.Sprintf("vapid t=%s,k=%s", signedJWT, b64Key))
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
	privKey := key.ToElliptic()
	return token.SignedString(privKey)
}

func VAPIDHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-control", "public, max-age=86400")

	keyPair := getOrCreateVAPIDKey()
	json.NewEncoder(w).Encode(&keyPair)
}

func getOrCreateVAPIDKey() *VAPIDKey {
	c := session.DB("dispotrains").C("pushKey")
	var keyPair VAPIDKey
	err := c.Find(nil).One(&keyPair)
	if err != nil {
		c.DropCollection()
		keyPair = createKeyPair()
		c.Insert(keyPair)
	}
	return &keyPair
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
	}
}
