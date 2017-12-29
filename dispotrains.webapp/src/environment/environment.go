package environment

import (
	"flag"
)

const (
	prodMongoDbAddress = "db"
	devMongoDbAddress  = "localhost"
)

var (
	prod        = flag.Bool("prod", false, "Is in production")
	initialized = false
)

func initialize() {
	flag.Parse()
	initialized = true
}

// GetMongoDbAddress returns the address of the MongoDB instance to use.
func GetMongoDbAddress() string {
	if !initialized {
		initialize()
	}
	if *prod {
		return prodMongoDbAddress
	}
	return devMongoDbAddress

}
