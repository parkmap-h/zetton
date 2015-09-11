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
	search(point GeoPoint, distance Meter) []Space
}

//-- DomainLogic --
type DomainContext struct {
	Err error
}

func (c *DomainContext) nearSpaces(point GeoPoint, searcher NearSpaceSearchService) []Space {
	distance := meter(100)
	return searcher.search(point, distance)
}
