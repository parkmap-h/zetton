package zetton

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"bytes"
	"encoding/json"
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

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Criticalf(format string, args ...interface{})
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
	memcacheKey := "nearspaces"
	c := self.App
	if item, err := memcache.Get(c, memcacheKey); err == memcache.ErrCacheMiss {
		return getSpacesOnCacheOrDatastore(c)
	} else if err != nil {
		return getSpacesOnCacheOrDatastore(c)
	} else {
		c.Debugf("use cache: ", memcacheKey)
		r := bytes.NewBuffer(item.Value)
		var spaces []InfraSpace
		err2 := json.NewDecoder(r).Decode(&spaces)
		return infraSpacesToSpaces(spaces), err2
	}
}

func infraSpacesToSpaces(spaces []InfraSpace) []Space {
	ret := make([]Space, len(spaces))
	for i, space := range spaces {
		ret[i] = Space(&space)
	}
	return ret
}

func getSpacesOnCacheOrDatastore(c appengine.Context) ([]Space, error) {
	memcacheKey := "nearspaces"
	c.Debugf("no hit cache: ", memcacheKey)
	spaces, err2 := getSpacesOnDatastore(c)
	var w bytes.Buffer
	if err := json.NewEncoder(&w).Encode(spaces); err != nil {
		return []Space{}, err
	}
	item := &memcache.Item{
		Key:   memcacheKey,
		Value: w.Bytes(),
	}
	if err2 != nil {
		return []Space{}, err2
	}
	if err3 := memcache.Set(c, item); err3 != nil {
		c.Errorf("error setting item: %v", err3)
		return []Space{}, err3
	}
	return spaces, nil
}

func getSpacesOnDatastore(c appengine.Context) ([]Space, error) {
	q := datastore.NewQuery("Space").Order("-CreateAt").Limit(100)
	var spaces []InfraSpace
	_, err := q.GetAll(c, &spaces)
	if err != nil {
		return nil, err
	}
	ret := infraSpacesToSpaces(spaces)
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
	prop := map[string]interface{}{spaceValueName: infra.Value_, "createAt": infra.CreateAt_.Unix()}
	return geojson.NewFeature(geom, prop, nil)
}

func SpacesToFeatureCollection(spaces []Space) *geojson.FeatureCollection {
	features := make([]*geojson.Feature, len(spaces))
	for i, space := range spaces {
		features[i] = SpaceToFeature(space)
	}
	return geojson.NewFeatureCollection(features)
}
