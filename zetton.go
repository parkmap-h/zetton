package zetton

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Space struct {
    Latitude float32 `json:"latitude"`
    Longitude float32 `json:"longitude"`
    Value int `json:"value"`
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
    dec := json.NewDecoder(r.Body)
    var space Space
    err := dec.Decode(&space)
    if err != nil {
        fmt.Fprint(w, "invalid json: " + err.Error())
        return
    }
    fmt.Fprint(w, space)
    return
}
