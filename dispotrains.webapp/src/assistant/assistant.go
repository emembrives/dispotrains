package assistant

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/emembrives/dispotrains/dispotrains.webapp/src/storage"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	//keyFile = "/dispotrains/key/dialogflow.json"
	keyFile = "../../Dispotrains-dialogflow.json"
)

var (
	saintRegexp = regexp.MustCompile("\\bst\\b")
)

func getJWTKey() (*http.Client, error) {
	d, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(d,
		"https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}
	return conf.Client(oauth2.NoContext), nil
}

type Entity struct {
	Value    string   `json:"value"`
	Synonyms []string `json:"synonyms"`
}

func makeSynonyms(name string) []string {
	synonyms := make([]string, 0)
	baseName := strings.ToLower(name)
	baseName = strings.TrimPrefix(baseName, "gare de ")
	baseName = strings.TrimPrefix(baseName, "gare d'")
	baseName = saintRegexp.ReplaceAllString(baseName, "saint")
	synonyms = append(synonyms, baseName)
	synonyms = append(synonyms, "gare de "+baseName)
	synonyms = append(synonyms, "gare d'"+baseName)
	synonyms = append(synonyms, "station de "+baseName)
	synonyms = append(synonyms, "station d'"+baseName)
	synonyms = append(synonyms, "station "+baseName)
	return synonyms
}

// UpdateStationList updates the station entity list in the Assistant.
func UpdateStationList(stations []*storage.Station) {
	entities := make([]Entity, len(stations))
	for i, station := range stations {
		entities[i].Value = station.Name
		entities[i].Synonyms = makeSynonyms(station.DisplayName)
	}

	client, err := getJWTKey()
	if err != nil {
		panic(err)
	}
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.Encode(map[string]interface{}{"entities": entities})
	log.Printf("Request: %s\n", b.String())
	resp, err := client.Post(
		"https://dialogflow.googleapis.com/v2beta1/projects/agent/entityTypes/Station/entities:batchUpdate", "application/json", &b)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("Response: %s\n", string(body))
}
