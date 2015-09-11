package zetton

import (
	"appengine"
	"github.com/kpawlik/geojson"
	"appengine/datastore"
)

type Latitude float64
type Longitude float64
type SpaceValueType int
type MeterValueType int

const spaceValueName = "value"

type InfraSpace struct {
	Point_ appengine.GeoPoint `datastore:"Point"`
	Value_ int `datastore:"Value,noindex"`
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
	app appengine.Context
	domain DomainContext
}

func (self *NearSpaceSearchServiceImpl) search(point GeoPoint, distance Meter) []Space {
	q := datastore.NewQuery("Space").Limit(100)
	var spaces []InfraSpace
	_, self.domain.Err = q.GetAll(self.app, &spaces)
	ret := make([]Space, len(spaces))
	for i, space := range spaces {
		ret[i] = Space(&space)
	}
	return ret
}

func featureToSpace(feature *geojson.Feature) *InfraSpace {
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
		ret[i] = *featureToSpace(feature)
	}
	return ret
}

func SpaceToFeature(space Space) *geojson.Feature {
	lng := geojson.CoordType(space.Point().Lng)
	lat := geojson.CoordType(space.Point().Lat)
	c := geojson.Coordinate{lng, lat}
	geom := geojson.NewPoint(c)
	prop := map[string]interface{}{spaceValueName: space.Value()}
	return geojson.NewFeature(geom, prop, nil)
}

func SpacesToFeatureCollection(spaces []Space) *geojson.FeatureCollection {
	features := make([]*geojson.Feature, len(spaces))
	for i, space := range spaces {
		features[i] = SpaceToFeature(space)
	}
	return geojson.NewFeatureCollection(features)
}
