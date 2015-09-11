package zetton

import (
	"encoding/json"
	"github.com/kpawlik/geojson"
	"net/http"
)

func listSpacesAction(ctx DomainContext, w http.ResponseWriter, r *http.Request) {
	point := GeoPoint{
		Lat: Latitude(132.3),
		Lng: Longitude(32.1),
	}

	spaces := ctx.queryNearSpace(point, ctx.NearSearchService)

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
	dec := json.NewDecoder(r.Body)
	var request geojson.Feature
	ctx.Err = dec.Decode(&request)
	if ctx.Err != nil {
		return
	}
	space := featureToInfraSpace(&request)
	ctx.commandCreateSpace(space)

	if ctx.Err != nil {
		return
	}
	err := json.NewEncoder(w).Encode(request)
	if err != nil {
		ctx.Err = InvalidJsonError{Err: err}
		return
	}
}
