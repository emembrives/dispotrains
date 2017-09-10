package push

import (
	"fmt"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"
)

const (
	pushRegistrationCollection = "pushSubscribers"
	pushKeyCollection          = "pushKey"
)

func Register(session *mgo.Session, subscription webpush.Subscription) {
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
	fmt.Printf("Sending to %d devices\n", len(registrations))
	for _, registration := range registrations {
		pushToOne(registration, vapid)
	}
}

func pushToOne(registration *storage.Registration, vapid *storage.VAPIDKey) {
	sub := registration.Subscription

	res, err := webpush.SendNotification([]byte("Test"), &sub, &webpush.Options{
		Subscriber:      "dispotrains@membrives.fr",
		VAPIDPrivateKey: vapid.PrivateKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
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
	public, private, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatal(err)
	}

	return storage.VAPIDKey{
		PublicKey:  public,
		PrivateKey: private,
	}
}
