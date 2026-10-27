package seeders

import (
	"fmt"
	"time"

	"leochat/app/facades"
	"leochat/app/models"
)

const (
	DemoEmail    = "demo@leochat.id"
	DemoPassword = "password123"
)

// Around Gedung Sate, Bandung.
const demoLat, demoLng = -6.9025, 107.6187

// DemoSeeder creates a demo account plus nearby profiles to swipe on. Some of
// them already liked the demo account, so liking them back creates a match.
type DemoSeeder struct{}

func (s *DemoSeeder) Signature() string { return "DemoSeeder" }

type demoProfile struct {
	name          string
	gender        models.Gender
	age           int
	bio           string
	dLat, dLng    float64
	likesDemoUser bool
}

var demoProfiles = []demoProfile{
	{"Rina", models.GenderFemale, 24, "Pecinta kopi dan hiking. Cari teman ngobrol yang seru.", 0.010, 0.012, true},
	{"Ayu", models.GenderFemale, 26, "Desainer grafis. Suka pameran seni dan kucing.", -0.020, 0.008, false},
	{"Nadia", models.GenderFemale, 23, "Mahasiswi farmasi, suka masak dan K-drama.", 0.030, -0.015, true},
	{"Sekar", models.GenderFemale, 28, "Fotografer lepas. Weekend biasanya ke Lembang.", -0.045, -0.030, false},
	{"Dewi", models.GenderFemale, 30, "Product manager. Lari pagi tiap Minggu di Car Free Day.", 0.060, 0.040, false},
	{"Putri", models.GenderFemale, 22, "Barista dan pemain ukulele amatir.", 0.005, -0.005, false},
	{"Bima", models.GenderMale, 27, "Software engineer. Main futsal tiap Kamis.", 0.015, 0.020, true},
	{"Raka", models.GenderMale, 25, "Pendaki gunung, sudah ke 12 puncak.", -0.012, -0.018, false},
	{"Dimas", models.GenderMale, 29, "Chef di restoran Sunda. Suka eksperimen sambal.", 0.040, -0.010, false},
	{"Arif", models.GenderMale, 31, "Dosen dan penulis. Koleksi buku bekas.", -0.035, 0.025, false},
}

func (s *DemoSeeder) Run() error {
	exists, err := facades.Orm().Query().Model(&models.User{}).Where("email = ?", DemoEmail).Exists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	password, err := facades.Hash().Make(DemoPassword)
	if err != nil {
		return err
	}

	demo := newUser("Demo Leo", DemoEmail, password, models.GenderMale, models.PreferenceEveryone, 27, demoLat, demoLng)
	demo.Bio = "Akun demo untuk mencoba aplikasi."
	if err := facades.Orm().Query().Create(&demo); err != nil {
		return err
	}

	for i, p := range demoProfiles {
		email := fmt.Sprintf("demo%02d@leochat.id", i+1)
		user := newUser(p.name, email, password, p.gender, models.PreferenceEveryone, p.age, demoLat+p.dLat, demoLng+p.dLng)
		user.Bio = p.bio
		if err := facades.Orm().Query().Create(&user); err != nil {
			return err
		}

		if p.likesDemoUser {
			like := models.Swipe{SwiperID: user.ID, SwipedID: demo.ID, Direction: models.SwipeLike}
			if err := facades.Orm().Query().Create(&like); err != nil {
				return err
			}
		}
	}
	return nil
}

func newUser(name, email, password string, gender models.Gender, pref models.Preference, age int, lat, lng float64) models.User {
	now := time.Now().UTC()
	birth := time.Date(now.Year()-age, time.March, 15, 0, 0, 0, 0, time.UTC)
	if birth.After(now.AddDate(-age, 0, 0)) {
		birth = birth.AddDate(-1, 0, 0)
	}

	return models.User{
		Name:         name,
		Email:        email,
		Password:     password,
		BirthDate:    birth,
		Gender:       gender,
		InterestedIn: pref,
		Latitude:     &lat,
		Longitude:    &lng,
		LastActiveAt: &now,
	}
}
