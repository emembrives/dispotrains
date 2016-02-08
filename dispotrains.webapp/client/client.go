package client

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/emembrives/tinkerings/dispotrains.webapp/storage"
	"golang.org/x/net/html"
)

func GetAllLines() ([]*storage.Line, error) {
	resp, err := http.Get("http://www.infomobi.com/fr/voyageurs-en-fauteuil/transports-publics-accessibles/gares-et-stations-accessibles/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyParser, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	var form *html.Node = findNode(bodyParser, "form")
	var optgroup *html.Node = findNode(form, "optgroup")
	var lines []*storage.Line = parseOptGroup(optgroup, optgroup.Attr[0].Val)
	optgroup = findNext(optgroup)
	lines = append(lines, parseOptGroup(optgroup, optgroup.Attr[0].Val)...)
	optgroup = findNext(optgroup)
	lines = append(lines, parseOptGroup(optgroup, optgroup.Attr[0].Val)...)
	for _, line := range lines {
		getStations(line)
	}
	return lines, nil
}

func parseOptGroup(optgroup *html.Node, network string) []*storage.Line {
	var lines []*storage.Line = make([]*storage.Line, 0)
	for option := optgroup.FirstChild; option != nil; option = option.NextSibling {
		if option.Type != html.ElementNode {
			continue
		}
		var code string = option.Attr[0].Val
		var name string = option.FirstChild.Data
		lines = append(lines, storage.NewLine(network, name, code))
	}
	return lines
}

func getStations(line *storage.Line) ([]*storage.Station, error) {
	var url *url.URL = line.GetURL()
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyParser, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	var stations []*storage.Station = make([]*storage.Station, 0)
	var table *html.Node = findNode(bodyParser, "table")
	if table == nil {
		return stations, nil
	}
	var tbody *html.Node = findNode(table, "tbody")
	var row *html.Node
	for row = findNode(tbody, "tr"); row != nil; row = findNext(row) {
		var col *html.Node = findNode(row, "td")
		var fullName string = col.FirstChild.Data
		var nameCity []string = strings.Split(fullName, ",")
		var name string = nameCity[0]
		var city string
		if len(nameCity) > 1 {
			city = nameCity[1]
		}
		col = findNext(col)
		//findAttrByKey(findNode(col, "p"), "class")
		col = findNext(col)
		var a *html.Node = findNode(col, "a")
		var station *storage.Station
		if a != nil {
			stationUrl, err := url.Parse(findAttrByKey(a, "href").Val)
			if err != nil {
				return nil, err
			}
			var code string = stationUrl.Query().Get("tx_stifinfomobi_pi3[externalcode]")
			station = storage.NewStation(name, city, code)
		} else {
			station = storage.NewRampStation(name, city)
		}
		station.AttachLine(line)
		getElevatorsAndStatus(station)
		if station.LastUpdate.After(line.LastUpdate) {
			line.LastUpdate = station.LastUpdate
		}
		stations = append(stations, station)
	}
	return stations, nil
}

func getElevatorsAndStatus(station *storage.Station) error {
	if !station.HasElevators {
		return nil
	}
	var url *url.URL = station.GetURL()
	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyParser, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}
	var contents *html.Node = findNodeWithAttributes(bodyParser, "div", map[string]string{"id": "contentRight"})
	var p *html.Node = findNode(contents, "p")
	if !strings.HasPrefix(p.FirstChild.Data, "Situation en date du") {
		p = findNext(p)
	}
	if !strings.HasPrefix(p.FirstChild.Data, "Situation en date du") {
		p = findNext(p)
	}
	var date string = p.FirstChild.Data[21:]
	var table *html.Node = findNode(bodyParser, "table")
	var tbody *html.Node = findNode(table, "tbody")
	var row *html.Node
	for row = findNode(tbody, "tr"); row != nil; row = findNext(row) {
		var col *html.Node = findNode(row, "td")
		var code, situation, direction, status string
		if col.FirstChild != nil {
			code = col.FirstChild.Data
		}
		col = findNext(col)
		if col.FirstChild != nil {
			situation = col.FirstChild.Data
		}
		col = findNext(col)
		if col.FirstChild != nil {
			direction = col.FirstChild.Data
		}
		col = findNext(col)
		if col.FirstChild != nil && findNode(col, "span") != nil {
			status = findNode(col, "span").FirstChild.Data
		} else if col.FirstChild != nil {
			status = col.FirstChild.Data
		}
		var elevator *storage.Elevator = station.NewElevator(code, situation, direction)
		_, err = elevator.NewStatus(status, date)
		if err != nil {
			return err
		}
		if elevator.Status.LastUpdate.After(station.LastUpdate) {
			station.LastUpdate = elevator.Status.LastUpdate
		}
	}
	return nil
}
