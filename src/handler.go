package zetton

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"github.com/kpawlik/geojson"
	"net/http"
)

func spacesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	var err error
	switch r.Method {
	case "GET":
		err = listSpacesHandler(w, r)
	case "POST":
		err = createSpacesHandler(w, r)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	}
	// unmatched Route
}

func listSpacesHandler(w http.ResponseWriter, r *http.Request) error {
	var err error
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Space").Limit(100)
	var spaces []Space
	_, err = q.GetAll(c, &spaces)
	if err != nil {
		return err
	}
	featureCollection := SpacesToFeatureCollection(spaces)
	err = json.NewEncoder(w).Encode(featureCollection)
	if err != nil {
		return InvalidJsonError{Err: err}
	}
	return nil
}

func createSpacesHandler(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	var err error
	dec := json.NewDecoder(r.Body)
	var request geojson.Feature
	err = dec.Decode(&request)
	if err != nil {
		return err
	}
	space := featureToSpace(&request)
	fmt.Fprint(w, space.Point.Lat)
	key := datastore.NewIncompleteKey(c, "Space", nil)
	_, err = datastore.Put(c, key, space)
	if err != nil {
		return err
	}
	fmt.Fprint(w, space)
	return nil
}
