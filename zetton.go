package zetton

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"github.com/kpawlik/geojson"
	"net/http"
)

type Space struct {
	Point appengine.GeoPoint
	Value int `datastore:",noindex"`
}

const spaceValueName = "value"

func featureToSpace(feature *geojson.Feature) *Space {
	point, _ := feature.GetGeometry()
	coordinate := point.(*geojson.Point).Coordinates
	prop := feature.Properties
	return &Space{
		Point: appengine.GeoPoint{
			Lng: float64(coordinate[0]),
			Lat: float64(coordinate[1]),
		},
		Value: int(prop[spaceValueName].(float64)),
	}
}

func FeatureColloctionToSpaces(featureCollection *geojson.FeatureCollection) []Space {
	features := featureCollection.Features
	ret := make([]Space, len(features))
	for i, feature := range features {
		ret[i] = *featureToSpace(feature)
	}
	return ret
}

func SpaceToFeature(space *Space) *geojson.Feature {
	lng := geojson.CoordType(space.Point.Lng)
	lat := geojson.CoordType(space.Point.Lat)
	c := geojson.Coordinate{lng, lat}
	geom := geojson.NewPoint(c)
	prop := map[string]interface{}{spaceValueName: space.Value}
	return geojson.NewFeature(geom, prop, nil)
}

func SpacesToFeatureCollection(spaces []Space) *geojson.FeatureCollection {
	features := make([]*geojson.Feature, len(spaces))
	for i, space := range spaces {
		features[i] = SpaceToFeature(&space)
	}
	return geojson.NewFeatureCollection(features)
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
	featureCollection := SpacesToFeatureCollection(spaces)
	err = json.NewEncoder(w).Encode(featureCollection)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	return
}

func createSpacesHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	dec := json.NewDecoder(r.Body)
	var request geojson.Feature
	err := dec.Decode(&request)
	if err != nil {
		fmt.Fprint(w, "invalid json: "+err.Error())
		return
	}
	space := featureToSpace(&request)
	fmt.Fprint(w, space.Point.Lat)
	key := datastore.NewIncompleteKey(c, "Space", nil)
	_, err2 := datastore.Put(c, key, space)
	if err2 != nil {
		fmt.Fprint(w, err2.Error())
		return
	}
	fmt.Fprint(w, space)
	return
}
