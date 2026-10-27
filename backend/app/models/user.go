package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type Preference string

const (
	PreferenceMale     Preference = "male"
	PreferenceFemale   Preference = "female"
	PreferenceEveryone Preference = "everyone"
)

// Accepts reports whether someone of gender g fits this preference.
func (p Preference) Accepts(g Gender) bool {
	return p == PreferenceEveryone || string(p) == string(g)
}

// Genders lists the genders this preference accepts.
func (p Preference) Genders() []Gender {
	if p == PreferenceEveryone {
		return []Gender{GenderMale, GenderFemale}
	}
	return []Gender{Gender(p)}
}

type User struct {
	orm.Model
	Name         string
	Email        string
	Password     string
	BirthDate    time.Time
	Gender       Gender
	InterestedIn Preference
	Bio          string
	Latitude     *float64
	Longitude    *float64
	LastActiveAt *time.Time
	Photos       []UserPhoto `gorm:"foreignKey:UserID"`
}

// AgeAt returns the user's age in whole years on the given date.
func (u *User) AgeAt(now time.Time) int {
	return YearsBetween(u.BirthDate, now)
}

func (u *User) HasLocation() bool {
	return u.Latitude != nil && u.Longitude != nil
}

// YearsBetween counts the full years elapsed from birth to now.
func YearsBetween(birth, now time.Time) int {
	years := now.Year() - birth.Year()
	if now.Month() < birth.Month() || (now.Month() == birth.Month() && now.Day() < birth.Day()) {
		years--
	}
	return years
}
