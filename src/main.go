package zetton

import (
	"_vendor/src/github.com/gorilla/mux"
	"appengine"
	"fmt"
	"net/http"
	"time"
)

const location = "Asia/Tokyo"

func init() {
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	r := mux.NewRouter()
	r.HandleFunc("/spaces/{year:20[0-9][0-9]}/{month:[0-1][0-9]}/{day:[0-3][0-9]}", createHandler(listSpacesByDayAction)).Methods("GET")
	r.HandleFunc("/spaces", createHandler(listSpacesAction)).Methods("GET")
	r.HandleFunc("/spaces", createHandler(createSpacesAction)).Methods("POST")
	http.Handle("/", r)
}

func createHandler(f func(w http.ResponseWriter, ctx *DomainContext)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		ctx := DomainContext{
			Request:           r,
			Logger:            c,
			SpaceRepository:   &SpaceRepositoryOnDatastore{C: c},
			NearSearchService: &NearSpaceSearchServiceImpl{App: c},
		}
		f(w, &ctx)
		if ctx.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, ctx.Err.Error())
		}
	}
}
