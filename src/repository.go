package zetton

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"time"
)

type SpaceKey int64

type SpaceRepository interface {
	store(key *SpaceKey, space Space) (SpaceKey, Space, error)
	resolveByDate(start time.Time) ([]Space, error)
}

type SpaceRepositoryOnDatastore struct {
	C appengine.Context
}

func (self *SpaceRepositoryOnDatastore) store(key *SpaceKey, space Space) (SpaceKey, Space, error) {
	var datastoreKey *datastore.Key
	infraSpace := spaceToInfra(space)
	if key != nil {
		datastore.NewKey(self.C, "Space", "", int64(*key), nil)
	} else {
		datastoreKey = datastore.NewIncompleteKey(self.C, "Space", nil)
		infraSpace.CreateAt_ = time.Now()
	}
	returnKey, err := datastore.Put(self.C, datastoreKey, infraSpace)
	memcacheKey := "nearspaces"
	memcache.Delete(self.C, memcacheKey)
	return SpaceKey(returnKey.IntID()), space, err
}

func (self *SpaceRepositoryOnDatastore) resolveByDate(start time.Time) ([]Space, error) {
	end := start.AddDate(0, 0, 1)
	q := datastore.NewQuery("Space").Order("-CreateAt").Filter("CreateAt >=", start).Filter("CreateAt <", end)
	var spaces []InfraSpace
	_, err := q.GetAll(self.C, &spaces)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	ret := make([]Space, len(spaces))
	for i, _ := range spaces {
		ret[i] = Space(&(spaces[i]))
	}
	return ret, err
}
