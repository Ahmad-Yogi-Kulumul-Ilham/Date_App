package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"leochat/app/http/middleware"
	"leochat/app/http/requests"
	"leochat/app/http/responses"
	"leochat/app/services"
)

type SwipeController struct {
	swipes *services.SwipeService
}

func NewSwipeController(swipes *services.SwipeService) *SwipeController {
	return &SwipeController{swipes: swipes}
}

type swipeResponse struct {
	Matched bool  `json:"matched"`
	MatchID *uint `json:"match_id"`
}

func (c *SwipeController) Store(ctx http.Context) http.Response {
	var req requests.SwipeRequest
	if res, ok := validate(ctx, &req); !ok {
		return res
	}

	result, err := c.swipes.Swipe(ctx, middleware.CurrentUserID(ctx), req.Target(), req.SwipeDirection())
	if err != nil {
		return responses.Error(ctx, err)
	}

	body := swipeResponse{Matched: result.Matched()}
	if result.Matched() {
		body.MatchID = &result.Match.ID
	}
	return responses.Created(ctx, body)
}
