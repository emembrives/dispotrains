package storage

import (
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
)

type VAPIDKey struct {
	PublicKey  string
	PrivateKey string
}

type Registration struct {
	Subscription webpush.Subscription
	LastUpdated  time.Time
}

type Subscription struct {
	Endpoint       string
	ExpirationTime interface{}
	Keys           SubscriptionKeys
}

type SubscriptionKeys struct {
	Auth   string
	P256dh string
}
