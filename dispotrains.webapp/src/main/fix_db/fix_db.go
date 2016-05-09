package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior
	// session.SetMode(mgo.Monotonic, true)

	c := session.DB("dispotrains").C("statuses")
	index := mgo.Index{
		Key:        []string{"state", "lastupdate", "elevator"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	
	_, err = c.RemoveAll(bson.M{"elevator": nil})
	if err != nil {
		panic(err)
	}
}

