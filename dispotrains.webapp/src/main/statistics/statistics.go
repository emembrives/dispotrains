package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mapFunc = `function() {
		var output = { "elevator": this.elevator,
		"state": this.state,
		"startdate": this.lastupdate,
		"enddate": this.lastupdate };
	emit(this.elevator, output);
}`
	reduceFunc = `function(key, values) {
	var results = [];
	var last_value = values[0];
	for (var value of values) {
		last_value.enddate = value.startdate;
		if (value.state != last_value.state) {
			results.push(last_value);
			last_value = value;
		}
	}
	results.push(last_value);
	var obj = {"elevator": key,
		"states": results
	};
	return obj;
}
`
)

const (
	server = "localhost"
)

func main() {
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cStatuses := session.DB("dispotrains").C("statuses")
	err = cStatuses.EnsureIndexKey("lastupdate")
	if err != nil {
		panic(err)
	}
	mapReduce := mgo.MapReduce{
		Map:    mapFunc,
		Reduce: reduceFunc,
		Out:    bson.M{"replace": "statistics"},
	}
	_, err = cStatuses.Find(nil).Sort("lastupdate").MapReduce(&mapReduce, nil)
	if err != nil {
		panic(err)
	}

	cStatistics := session.DB("dispotrains").C("statistics")
	err = cStatistics.EnsureIndexKey("elevator")
	if err != nil {
		panic(err)
	}
}
