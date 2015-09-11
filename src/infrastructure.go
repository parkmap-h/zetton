package zetton

import (
	"appengine"
	"github.com/kpawlik/geojson"
)

const spaceValueName = "value"

func featureToSpace(feature *geojson.Feature) *Space {
	point, _ := feature.GetGeometry()
	coordinate := point.(*geojson.Point).Coordinates
	prop := feature.Properties
	return &Space{
		Point: appengine.GeoPoint{
			Lng: float64(coordinate[0]),
			Lat: float64(coordinate[1]),
		},
		Value: int(prop[spaceValueName].(float64)),
	}
}

func FeatureColloctionToSpaces(featureCollection *geojson.FeatureCollection) []Space {
	features := featureCollection.Features
	ret := make([]Space, len(features))
	for i, feature := range features {
		ret[i] = *featureToSpace(feature)
	}
	return ret
}

func SpaceToFeature(space *Space) *geojson.Feature {
	lng := geojson.CoordType(space.Point.Lng)
	lat := geojson.CoordType(space.Point.Lat)
	c := geojson.Coordinate{lng, lat}
	geom := geojson.NewPoint(c)
	prop := map[string]interface{}{spaceValueName: space.Value}
	return geojson.NewFeature(geom, prop, nil)
}

func SpacesToFeatureCollection(spaces []Space) *geojson.FeatureCollection {
	features := make([]*geojson.Feature, len(spaces))
	for i, space := range spaces {
		features[i] = SpaceToFeature(&space)
	}
	return geojson.NewFeatureCollection(features)
}
