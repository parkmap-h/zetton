package zetton

import (
	"net/http"
	"fmt"
)

func init() {
	http.HandleFunc("/spaces", spacesHandler)
}

func spacesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	ctx := DomainContext{}
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
