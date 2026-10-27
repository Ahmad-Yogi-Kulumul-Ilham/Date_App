package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"leochat/app/http/middleware"
	"leochat/app/http/requests"
	"leochat/app/http/resources"
	"leochat/app/http/responses"
	"leochat/app/services"
)

const photoDisk = "public"

type ProfileController struct {
	profiles  *services.ProfileService
	presenter *resources.Presenter
}

func NewProfileController(profiles *services.ProfileService, presenter *resources.Presenter) *ProfileController {
	return &ProfileController{profiles: profiles, presenter: presenter}
}

func (c *ProfileController) Show(ctx http.Context) http.Response {
	user, err := c.profiles.Get(ctx, middleware.CurrentUserID(ctx))
	if err != nil {
		return responses.Error(ctx, err)
	}
	return responses.OK(ctx, c.presenter.Me(user))
}

func (c *ProfileController) Update(ctx http.Context) http.Response {
	var req requests.UpdateProfileRequest
	if res, ok := validate(ctx, &req); !ok {
		return res
	}

	user, err := c.profiles.Update(ctx, middleware.CurrentUserID(ctx), req.Input())
	if err != nil {
		return responses.Error(ctx, err)
	}
	return responses.OK(ctx, c.presenter.Me(user))
}

func (c *ProfileController) UpdateLocation(ctx http.Context) http.Response {
	var req requests.UpdateLocationRequest
	if res, ok := validate(ctx, &req); !ok {
		return res
	}

	if err := c.profiles.UpdateLocation(ctx, middleware.CurrentUserID(ctx), req.Point()); err != nil {
		return responses.Error(ctx, err)
	}
	return responses.NoContent(ctx)
}

func (c *ProfileController) UploadPhoto(ctx http.Context) http.Response {
	var req requests.UploadPhotoRequest
	if res, ok := validate(ctx, &req); !ok {
		return res
	}

	file, err := ctx.Request().File("photo")
	if err != nil {
		return responses.BadRequest(ctx)
	}

	photo, err := c.profiles.AddPhoto(ctx, middleware.CurrentUserID(ctx), file.Disk(photoDisk))
	if err != nil {
		return responses.Error(ctx, err)
	}
	return responses.Created(ctx, c.presenter.Photo(*photo))
}

func (c *ProfileController) DeletePhoto(ctx http.Context) http.Response {
	photoID, ok := routeID(ctx, "id")
	if !ok {
		return notFound(ctx)
	}

	if err := c.profiles.DeletePhoto(ctx, middleware.CurrentUserID(ctx), photoID); err != nil {
		return responses.Error(ctx, err)
	}
	return responses.NoContent(ctx)
}
