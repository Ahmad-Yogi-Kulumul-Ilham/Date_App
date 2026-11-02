package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"leochat/app/facades"
	"leochat/app/http/middleware"
	"leochat/app/http/requests"
	"leochat/app/http/responses"
	"leochat/app/services"
)

type AccountController struct {
	accounts *services.AccountService
}

func NewAccountController(accounts *services.AccountService) *AccountController {
	return &AccountController{accounts: accounts}
}

// Destroy permanently deletes the signed-in user's account and revokes the
// token used for the request.
func (c *AccountController) Destroy(ctx http.Context) http.Response {
	var req requests.DeleteAccountRequest
	if res, ok := validate(ctx, &req); !ok {
		return res
	}

	if err := c.accounts.Delete(ctx, middleware.CurrentUserID(ctx), req.Password); err != nil {
		return responses.Error(ctx, err)
	}

	// The account is gone either way, so a failed revoke is logged, not returned.
	if err := facades.Auth(ctx).Logout(); err != nil {
		facades.Log().WithContext(ctx).Error(err)
	}
	return responses.NoContent(ctx)
}
