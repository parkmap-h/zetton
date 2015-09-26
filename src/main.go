package zetton

import (
	"appengine"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/spaces", spacesHandler)
}

func spacesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c := appengine.NewContext(r)
	ctx := DomainContext{
		SpaceRepository:   &SpaceRepositoryOnDatastore{C: c},
		NearSearchService: &NearSpaceSearchServiceImpl{App: c},
	}
	switch r.Method {
	case "GET":
		listSpacesAction(ctx, w, r)
	case "POST":
		createSpacesAction(ctx, w, r)
	}

	if ctx.Err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ctx.Err.Error())
	}
	// unmatched Route
}
