package zetton

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"github.com/kpawlik/geojson"
	"net/http"
)

func listSpacesAction(ctx DomainContext, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	point := GeoPoint{
		Lat: Latitude(132.3),
		Lng: Longitude(32.1),
	}
	searcher := NearSpaceSearchServiceImpl{app: c}
	spaces := ctx.nearSpaces(point, &searcher)
	if ctx.Err != nil {
		return
	}
	featureCollection := SpacesToFeatureCollection(spaces)
	err := json.NewEncoder(w).Encode(featureCollection)
	if err != nil {
		ctx.Err = InvalidJsonError{Err: err}
		return
	}
}

func createSpacesAction(ctx DomainContext, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	dec := json.NewDecoder(r.Body)
	var request geojson.Feature
	ctx.Err = dec.Decode(&request)
	if ctx.Err != nil {
		return
	}
	space := featureToSpace(&request)
	key := datastore.NewIncompleteKey(c, "Space", nil)
	_, ctx.Err = datastore.Put(c, key, space)
	if ctx.Err != nil {
		return
	}
	fmt.Fprint(w, space)
}
