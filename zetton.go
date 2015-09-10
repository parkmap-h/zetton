package zetton

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"net/http"
)

type SpaceJson struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Value     int     `json:"value"`
}

type Space struct {
	Point appengine.GeoPoint
	Value int `datastore:",noindex"`
}

func jsonToSpace(request *SpaceJson, space *Space) {
	space.Point = appengine.GeoPoint{Lat: request.Latitude, Lng: request.Longitude}
	space.Value = request.Value
}

func SpaceTojson(space *Space, request *SpaceJson) {
	request.Latitude = space.Point.Lat
	request.Longitude = space.Point.Lng
	request.Value = space.Value
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/spaces", spacesHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func spacesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listSpacesHandler(w, r)
	case "POST":
		createSpacesHandler(w, r)
		return
	}
	// unmatched Route
}

func listSpacesHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Space").Limit(100)
	var spaces []Space
	_, err := q.GetAll(c, &spaces)
	if err != nil {
		fmt.Fprint(w, "fail get spaces")
		return
	}
	var response = make([]SpaceJson, len(spaces))
	for i, space := range spaces {
		SpaceTojson(&space, &response[i])
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	return
}

func createSpacesHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	dec := json.NewDecoder(r.Body)
	var request SpaceJson
	err := dec.Decode(&request)
	if err != nil {
		fmt.Fprint(w, "invalid json: "+err.Error())
		return
	}
	var space Space
	jsonToSpace(&request, &space)
	key := datastore.NewIncompleteKey(c, "Space", nil)
	_, err2 := datastore.Put(c, key, &space)
	if err2 != nil {
		fmt.Fprint(w, err2.Error())
		return
	}
	fmt.Fprint(w, space)
	return
}
