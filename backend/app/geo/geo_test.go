package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoundingBoxAround(t *testing.T) {
	bandung := Point{Lat: -6.9025, Lng: 107.6187}

	box := BoundingBoxAround(bandung, 10)

	// One degree of latitude is ~111 km, so 10 km is ~0.09 degrees.
	assert.InDelta(t, 0.0898, bandung.Lat-box.MinLat, 0.001)
	assert.InDelta(t, 0.0898, box.MaxLat-bandung.Lat, 0.001)
	// Near the equator a degree of longitude is almost as long as one of latitude.
	assert.InDelta(t, 0.0905, box.MaxLng-bandung.Lng, 0.001)
}

func TestBoundingBoxAround_ClampsLatitudeAtPoles(t *testing.T) {
	box := BoundingBoxAround(Point{Lat: 89.99, Lng: 0}, 50)

	assert.Equal(t, 90.0, box.MaxLat)
	assert.Equal(t, 180.0, box.MaxLng)
}

func TestDisplayKm(t *testing.T) {
	tests := map[float64]int{0: 1, 400: 1, 1499: 1, 1500: 2, 12_345: 12}
	for meters, want := range tests {
		assert.Equal(t, want, DisplayKm(meters), "meters=%v", meters)
	}
}
