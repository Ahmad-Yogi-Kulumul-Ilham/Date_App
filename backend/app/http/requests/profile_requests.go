package requests

import (
	"github.com/goravel/framework/contracts/http"

	"leochat/app/geo"
	"leochat/app/models"
	"leochat/app/services"
)

// UpdateProfileRequest is a partial update: only fields present in the body change.
type UpdateProfileRequest struct {
	Name         *string `form:"name" json:"name"`
	Bio          *string `form:"bio" json:"bio"`
	InterestedIn *string `form:"interested_in" json:"interested_in"`
}

func (r *UpdateProfileRequest) Authorize(http.Context) error { return nil }

func (r *UpdateProfileRequest) Rules(http.Context) map[string]any {
	return map[string]any{
		"name":          "sometimes|required|string|min:2|max:50",
		"bio":           "sometimes|string|max:500",
		"interested_in": "sometimes|required|in:male,female,everyone",
	}
}

func (r *UpdateProfileRequest) Input() services.UpdateProfileInput {
	in := services.UpdateProfileInput{Name: r.Name, Bio: r.Bio}
	if r.InterestedIn != nil {
		pref := models.Preference(*r.InterestedIn)
		in.InterestedIn = &pref
	}
	return in
}

type UpdateLocationRequest struct {
	Latitude  float64 `form:"latitude" json:"latitude"`
	Longitude float64 `form:"longitude" json:"longitude"`
}

func (r *UpdateLocationRequest) Authorize(http.Context) error { return nil }

func (r *UpdateLocationRequest) Rules(http.Context) map[string]any {
	return map[string]any{
		"latitude":  "required|numeric|between:-90,90",
		"longitude": "required|numeric|between:-180,180",
	}
}

func (r *UpdateLocationRequest) Point() geo.Point {
	return geo.Point{Lat: r.Latitude, Lng: r.Longitude}
}

const maxPhotoKilobytes = 5 * 1024

type UploadPhotoRequest struct{}

func (r *UploadPhotoRequest) Authorize(http.Context) error { return nil }

func (r *UploadPhotoRequest) Rules(http.Context) map[string]any {
	return map[string]any{
		"photo": "required|file|mimes:jpg,jpeg,png,webp|max:" + itoa(maxPhotoKilobytes),
	}
}
