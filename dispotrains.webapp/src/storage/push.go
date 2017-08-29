package storage

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"math/big"
	"time"

	"gopkg.in/mgo.v2/bson"
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
	Subscription Subscription
	LastUpdated  time.Time
}

type Subscription struct {
	Endpoind       string
	ExpirationTime interface{}
	Keys           SubscriptionKeys
}

type SubscriptionKeys struct {
	Auth   string
	P256dh string
}
