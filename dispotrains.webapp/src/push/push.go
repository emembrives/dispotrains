package push

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/dgrijalva/jwt-go"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
)

const (
	pushRegistrationCollection = "pushSubscribers"
	pushKeyCollection          = "pushKey"
)

func Register(session *mgo.Session, subscription storage.Subscription) {
	registration := storage.Registration{
		Subscription: subscription,
		LastUpdated:  time.Now(),
	}
	c := session.DB("dispotrains").C(pushRegistrationCollection)
	_, err := c.Upsert(map[string]interface{}{"subscription": registration.Subscription}, registration)
	if err != nil {
		log.Println(err)
	}
}

func PushToAll(session *mgo.Session) {
	vapid := GetOrCreateVAPIDKey(session)

	c := session.DB("dispotrains").C(pushRegistrationCollection)
	registrations := make([]*storage.Registration, 0)
	c.Find(nil).All(&registrations)
	for _, registration := range registrations {
		pushToOne(registration, vapid)
	}
}

func pushToOne(registration *storage.Registration, vapid *storage.VAPIDKey) {
	signedJWT, err := createSignedJWT(vapid, "")
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

func createSignedJWT(vapidKey *storage.VAPIDKey, endpoint string) (string, error) {
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	endpointUrl.Path = ""
	endpointUrl.RawQuery = ""
	endpointUrl.Fragment = ""
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{
		Audience:  endpointUrl.String(),
		Subject:   "mailto:foo@bar.fr",
		ExpiresAt: time.Now().Add(time.Duration(12) * time.Hour).Unix(),
	})
	privKey := vapidKey.ToElliptic()
	return token.SignedString(privKey)
}

func GetOrCreateVAPIDKey(session *mgo.Session) *storage.VAPIDKey {
	c := session.DB("dispotrains").C(pushKeyCollection)
	var keyPair storage.VAPIDKey
	err := c.Find(nil).One(&keyPair)
	if err != nil {
		c.DropCollection()
		keyPair = createKeyPair()
		c.Insert(keyPair)
	}
	return &keyPair
}

func createKeyPair() storage.VAPIDKey {
	curve := elliptic.P256()
	// priv
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}
	// pub
	return storage.VAPIDKey{
		PrivateKeyX: priv.X,
		PrivateKeyY: priv.Y,
		PrivateKeyD: priv.D,
	}
}
