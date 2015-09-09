package zetton

import (
    "appengine"
    "encoding/json"
    "fmt"
    "net/http"
    "appengine/datastore"
)

type SpacePostRequest struct {
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Value int `json:"value"`
}

type Space struct {
    Point appengine.GeoPoint
    Value int `datastore:",noindex"`
}

func RequestToSpace(request *SpacePostRequest, space *Space) {
    space.Point = appengine.GeoPoint{Lat: request.Latitude, Lng: request.Longitude}
    space.Value = request.Value
}

func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/spaces", spacesHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}

func spacesHandler (w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        createSpacesHandler(w, r)
        return
    }
    // unmatched Route
}

func createSpacesHandler (w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    dec := json.NewDecoder(r.Body)
    var request SpacePostRequest
    err := dec.Decode(&request)
    if err != nil {
        fmt.Fprint(w, "invalid json: " + err.Error())
        return
    }
    var space Space
    RequestToSpace(&request, &space)
    key := datastore.NewIncompleteKey(c, "Space", nil)
    _, err2 := datastore.Put(c, key, &space)
    if err2 != nil {
        fmt.Fprint(w, err2.Error())
        return
    }
    fmt.Fprint(w, space)
    return
}
