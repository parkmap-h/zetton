package zetton

import (
	"_vendor/src/github.com/gorilla/mux"
	"encoding/json"
	"github.com/kpawlik/geojson"
	"net/http"
	"strconv"
	"time"
)

func listSpacesAction(w http.ResponseWriter, ctx *DomainContext) {
	ctx.Logger.Debugf("start listSpacesAction")
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
	ctx.Logger.Debugf("finish listSpacesAction")
}

func listSpacesByDayAction(w http.ResponseWriter, ctx *DomainContext) {
	vars := mux.Vars(ctx.Request)
	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	day, _ := strconv.Atoi(vars["day"])
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

	spaces, _ := ctx.SpaceRepository.resolveByDate(start)
	featureCollection := SpacesToFeatureCollection(spaces)
	err := json.NewEncoder(w).Encode(featureCollection)
	if err != nil {
		ctx.Err = InvalidJsonError{Err: err}
		return
	}
}

func createSpacesAction(w http.ResponseWriter, ctx *DomainContext) {
	dec := json.NewDecoder(ctx.Request.Body)
	var request geojson.Feature
	ctx.Err = dec.Decode(&request)
	if ctx.Err != nil {
		return
	}
	space := featureToInfraSpace(&request)
	ctx.commandCreateSpace(space)
	response := SpaceToFeature(space)

	if ctx.Err != nil {
		return
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Err = InvalidJsonError{Err: err}
		return
	}
}
