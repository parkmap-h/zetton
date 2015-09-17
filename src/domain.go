package zetton

//-- ValueObject --
type GeoPoint struct {
	Lat Latitude
	Lng Longitude
}

type Meter struct {
	Value MeterValueType
}

//-- Entity --
const SpaceValueMin = SpaceValueType(0)

type Space interface {
	Point() GeoPoint
	Value() SpaceValueType
}

// Factory
func meter(n int) Meter {
	return Meter{Value: MeterValueType(n)}
}

//-- DomainService --
type NearSpaceSearchService interface {
	search(point GeoPoint, distance Meter) ([]Space, error)
}

//-- DomainLogic --
type DomainContext struct {
	Err               error
	SpaceRepository   SpaceRepository
	NearSearchService NearSpaceSearchService
}

func (c *DomainContext) commandCreateSpace(space Space) Space {
	min := SpaceValueMin
	if int(min) > int(space.Value()) {
		c.Err = SpaceValueError{Message: "空き数は0以上の値を設定してください"}
		return nil
	}
	_, ret, err := c.SpaceRepository.store(nil, space)
	c.Err = err
	return ret
}

func (c *DomainContext) queryNearSpace(point GeoPoint, searcher NearSpaceSearchService) []Space {
	distance := meter(100)
	spaces, err := searcher.search(point, distance)
	c.Err = err
	return spaces
}
