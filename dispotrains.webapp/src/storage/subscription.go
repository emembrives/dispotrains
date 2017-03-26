package storage

import (
	"math/big"
	"time"
)

type Subscription struct {
	LastUpdate time.Time
	PrivateKey []byte
	x          *big.Int
	y          *big.Int
}
