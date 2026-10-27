package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"leochat/app/http/responses"
)

// validate binds and validates req. When it returns false, the caller must
// return the accompanying response unchanged.
func validate(ctx http.Context, req http.FormRequest) (http.Response, bool) {
	errs, err := ctx.Request().ValidateRequest(req)
	if err != nil {
		return responses.BadRequest(ctx), false
	}
	if errs != nil {
		return responses.ValidationFailed(ctx, errs), false
	}
	return nil, true
}

// routeID reads a positive numeric route parameter; ok is false for anything else.
func routeID(ctx http.Context, name string) (uint, bool) {
	id := ctx.Request().RouteInt64(name)
	if id <= 0 {
		return 0, false
	}
	return uint(id), true
}

func notFound(ctx http.Context) http.Response {
	return responses.Fail(ctx, http.StatusNotFound, "not_found", "Data tidak ditemukan.")
}
