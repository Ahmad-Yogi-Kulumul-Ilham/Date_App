package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"leochat/app/http/middleware"
	"leochat/app/http/requests"
	"leochat/app/http/responses"
	"leochat/app/services"
)

type SafetyController struct {
	safety *services.SafetyService
}

func NewSafetyController(safety *services.SafetyService) *SafetyController {
	return &SafetyController{safety: safety}
}

// Block hides both users from each other and ends their match. Repeating it is harmless.
func (c *SafetyController) Block(ctx http.Context) http.Response {
	var req requests.BlockRequest
	if res, ok := validate(ctx, &req); !ok {
		return res
	}

	if err := c.safety.Block(ctx, middleware.CurrentUserID(ctx), req.Target()); err != nil {
		return responses.Error(ctx, err)
	}
	return responses.NoContent(ctx)
}
