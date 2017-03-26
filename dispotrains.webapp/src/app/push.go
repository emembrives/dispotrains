package main

import (
	"crypto/elliptic"
	"crypto/rand"
)

func createSubscription() {
	curve := elliptic.P256()
	priv, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	pub := elliptic.Marshal(curve, x, y)
}
