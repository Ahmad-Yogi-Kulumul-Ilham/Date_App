// Package geo holds the small amount of geometry the discovery feature needs.
package geo

import "math"

const kmPerDegreeLatitude = 111.32

type Point struct {
	Lat float64
	Lng float64
}

type BoundingBox struct {
	MinLat, MaxLat float64
	MinLng, MaxLng float64
}

// BoundingBoxAround returns a box that fully contains the circle of radiusKm
// around center. It is used as a cheap, index-friendly prefilter before the
// exact spherical distance is computed in SQL.
func BoundingBoxAround(center Point, radiusKm float64) BoundingBox {
	latDelta := radiusKm / kmPerDegreeLatitude

	cosLat := math.Cos(center.Lat * math.Pi / 180)
	// Near the poles cos(lat) approaches zero; fall back to the full longitude range.
	lngDelta := 180.0
	if cosLat > 0.01 {
		lngDelta = math.Min(180, radiusKm/(kmPerDegreeLatitude*cosLat))
	}

	return BoundingBox{
		MinLat: math.Max(-90, center.Lat-latDelta),
		MaxLat: math.Min(90, center.Lat+latDelta),
		MinLng: center.Lng - lngDelta,
		MaxLng: center.Lng + lngDelta,
	}
}

// DisplayKm rounds a distance for display. Exact distances are never exposed so
// that a user's position cannot be triangulated; anything under 1 km shows as 1.
func DisplayKm(meters float64) int {
	km := int(math.Round(meters / 1000))
	if km < 1 {
		return 1
	}
	return km
}
