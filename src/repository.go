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
