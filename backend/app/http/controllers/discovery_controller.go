package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"leochat/app/http/middleware"
	"leochat/app/http/resources"
	"leochat/app/http/responses"
	"leochat/app/services"
)

type DiscoveryController struct {
	discovery *services.DiscoveryService
	presenter *resources.Presenter
}

func NewDiscoveryController(discovery *services.DiscoveryService, presenter *resources.Presenter) *DiscoveryController {
	return &DiscoveryController{discovery: discovery, presenter: presenter}
}

// Index lists nearby candidates. Query parameters are optional and clamped by
// the service: radius_km, min_age, max_age, limit.
func (c *DiscoveryController) Index(ctx http.Context) http.Response {
	req := ctx.Request()
	filter := services.DiscoverFilter{
		RadiusKm: req.QueryInt("radius_km"),
		MinAge:   req.QueryInt("min_age"),
		MaxAge:   req.QueryInt("max_age"),
		Limit:    req.QueryInt("limit"),
	}

	candidates, err := c.discovery.Discover(ctx, middleware.CurrentUserID(ctx), filter)
	if err != nil {
		return responses.Error(ctx, err)
	}
	return responses.OK(ctx, c.presenter.Candidates(candidates))
}
