package zetton

import (
	"appengine"
	"appengine/datastore"
	"github.com/kpawlik/geojson"
	"time"
)

type Latitude float64
type Longitude float64
type SpaceValueType int
type MeterValueType int

const spaceValueName = "value"

type InfraSpace struct {
	Point_    appengine.GeoPoint `datastore:"Point"`
	Value_    int                `datastore:"Value,noindex"`
	CreateAt_ time.Time          `datastore:"CreateAt"`
}

func (space *InfraSpace) Point() GeoPoint {
	return GeoPoint{
		Lat: Latitude(space.Point_.Lat),
		Lng: Longitude(space.Point_.Lng),
	}
}

func (space *InfraSpace) Value() SpaceValueType {
	return SpaceValueType(space.Value_)
}

type NearSpaceSearchServiceImpl struct {
	App appengine.Context
}

func (self *NearSpaceSearchServiceImpl) search(point GeoPoint, distance Meter) ([]Space, error) {
	q := datastore.NewQuery("Space").Order("-CreateAt").Limit(100)
	var spaces []InfraSpace
	_, err := q.GetAll(self.App, &spaces)
	if err != nil {
		return nil, err
	}
	ret := make([]Space, len(spaces))
	for i, _ := range spaces {
		ret[i] = Space(&(spaces[i]))
	}
	println(ret[0].Value())
	return ret, nil
}

func featureToInfraSpace(feature *geojson.Feature) *InfraSpace {
	point, _ := feature.GetGeometry()
	coordinate := point.(*geojson.Point).Coordinates
	prop := feature.Properties
	return &InfraSpace{
		Point_: appengine.GeoPoint{
			Lng: float64(coordinate[0]),
			Lat: float64(coordinate[1]),
		},
		Value_: int(prop[spaceValueName].(float64)),
	}
}

func FeatureColloctionToSpaces(featureCollection *geojson.FeatureCollection) []InfraSpace {
	features := featureCollection.Features
	ret := make([]InfraSpace, len(features))
	for i, feature := range features {
		ret[i] = *featureToInfraSpace(feature)
	}
	return ret
}

func spaceToInfra(space Space) *InfraSpace {
	switch t := space.(type) {
	case *InfraSpace:
		return t
	}
	return nil
}

func SpaceToFeature(space Space) *geojson.Feature {
	infra := spaceToInfra(space)
	lng := geojson.CoordType(infra.Point_.Lng)
	lat := geojson.CoordType(infra.Point_.Lat)
	c := geojson.Coordinate{lng, lat}
	geom := geojson.NewPoint(c)
	prop := map[string]interface{}{spaceValueName: infra.Value_, "createAt": infra.CreateAt_}
	return geojson.NewFeature(geom, prop, nil)
}

func SpacesToFeatureCollection(spaces []Space) *geojson.FeatureCollection {
	features := make([]*geojson.Feature, len(spaces))
	for i, space := range spaces {
		println(spaces)
		features[i] = SpaceToFeature(space)
	}
	return geojson.NewFeatureCollection(features)
}
