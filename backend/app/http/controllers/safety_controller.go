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

// Report flags a user for moderators and blocks them.
func (c *SafetyController) Report(ctx http.Context) http.Response {
	var req requests.ReportRequest
	if res, ok := validate(ctx, &req); !ok {
		return res
	}

	err := c.safety.Report(ctx, middleware.CurrentUserID(ctx), req.Target(), req.ReportReason(), req.Details)
	if err != nil {
		return responses.Error(ctx, err)
	}
	return responses.NoContent(ctx)
}
