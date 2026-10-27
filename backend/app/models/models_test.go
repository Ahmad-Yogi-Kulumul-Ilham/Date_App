package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestYearsBetween(t *testing.T) {
	birth := time.Date(2000, time.March, 15, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, 25, YearsBetween(birth, time.Date(2026, time.March, 14, 0, 0, 0, 0, time.UTC)))
	assert.Equal(t, 26, YearsBetween(birth, time.Date(2026, time.March, 15, 0, 0, 0, 0, time.UTC)))
}

func TestPreference(t *testing.T) {
	assert.True(t, PreferenceEveryone.Accepts(GenderMale))
	assert.True(t, PreferenceFemale.Accepts(GenderFemale))
	assert.False(t, PreferenceFemale.Accepts(GenderMale))
	assert.Equal(t, []Gender{GenderMale}, PreferenceMale.Genders())
}

func TestMatchParticipants(t *testing.T) {
	one, two := OrderedPair(9, 4)
	m := Match{UserOneID: one, UserTwoID: two}

	assert.Equal(t, uint(4), m.UserOneID)
	assert.True(t, m.Includes(9))
	assert.False(t, m.Includes(5))
	assert.Equal(t, uint(4), m.PartnerOf(9))
}
