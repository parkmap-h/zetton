package zetton

import "appengine"

type Space struct {
	Point appengine.GeoPoint
	Value int `datastore:",noindex"`
}
