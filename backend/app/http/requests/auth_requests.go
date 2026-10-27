package requests

import (
	"time"

	"github.com/goravel/framework/contracts/http"

	"leochat/app/models"
	"leochat/app/services"
)

const dateLayout = "2006-01-02"

type RegisterRequest struct {
	Name         string `form:"name" json:"name"`
	Email        string `form:"email" json:"email"`
	Password     string `form:"password" json:"password"`
	BirthDate    string `form:"birth_date" json:"birth_date"`
	Gender       string `form:"gender" json:"gender"`
	InterestedIn string `form:"interested_in" json:"interested_in"`
}

func (r *RegisterRequest) Authorize(http.Context) error { return nil }

func (r *RegisterRequest) Rules(http.Context) map[string]any {
	return map[string]any{
		"name":  "required|string|min:2|max:50",
		"email": "required|email|max:255",
		// bcrypt only uses the first 72 bytes of a password.
		"password":      "required|string|min:8|max:72",
		"birth_date":    "required|date_format:" + dateLayout,
		"gender":        "required|in:male,female",
		"interested_in": "required|in:male,female,everyone",
	}
}

func (r *RegisterRequest) Input() (services.RegisterInput, error) {
	birthDate, err := time.Parse(dateLayout, r.BirthDate)
	if err != nil {
		return services.RegisterInput{}, err
	}

	return services.RegisterInput{
		Name:         r.Name,
		Email:        r.Email,
		Password:     r.Password,
		BirthDate:    birthDate,
		Gender:       models.Gender(r.Gender),
		InterestedIn: models.Preference(r.InterestedIn),
	}, nil
}

type LoginRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (r *LoginRequest) Authorize(http.Context) error { return nil }

func (r *LoginRequest) Rules(http.Context) map[string]any {
	return map[string]any{
		"email":    "required|email",
		"password": "required|string",
	}
}
