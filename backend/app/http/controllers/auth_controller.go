package controllers

import (
	"errors"

	"github.com/goravel/framework/contracts/http"
	frameworkerrors "github.com/goravel/framework/errors"

	"leochat/app/facades"
	"leochat/app/http/requests"
	"leochat/app/http/resources"
	"leochat/app/http/responses"
	"leochat/app/models"
	"leochat/app/services"
)

type AuthController struct {
	auth      *services.AuthService
	presenter *resources.Presenter
}

func NewAuthController(auth *services.AuthService, presenter *resources.Presenter) *AuthController {
	return &AuthController{auth: auth, presenter: presenter}
}

func (c *AuthController) Register(ctx http.Context) http.Response {
	var req requests.RegisterRequest
	if res, ok := validate(ctx, &req); !ok {
		return res
	}

	input, err := req.Input()
	if err != nil {
		return responses.BadRequest(ctx)
	}

	user, err := c.auth.Register(ctx, input)
	if err != nil {
		return responses.Error(ctx, err)
	}

	token, err := c.issueToken(ctx, user)
	if err != nil {
		return responses.Error(ctx, err)
	}
	return responses.Created(ctx, token)
}

func (c *AuthController) Login(ctx http.Context) http.Response {
	var req requests.LoginRequest
	if res, ok := validate(ctx, &req); !ok {
		return res
	}

	user, err := c.auth.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return responses.Error(ctx, err)
	}

	token, err := c.issueToken(ctx, user)
	if err != nil {
		return responses.Error(ctx, err)
	}
	return responses.OK(ctx, token)
}

// Refresh accepts an expired (but still refreshable) token and returns a new one.
func (c *AuthController) Refresh(ctx http.Context) http.Response {
	guard := facades.Auth(ctx)
	if _, err := guard.Parse(ctx.Request().Header("Authorization")); err != nil && !isExpired(err) {
		return responses.Fail(ctx, http.StatusUnauthorized, "unauthenticated", "Token tidak valid.")
	}

	token, err := guard.Refresh()
	if err != nil {
		return responses.Fail(ctx, http.StatusUnauthorized, "refresh_expired", "Sesi sudah berakhir, silakan login ulang.")
	}
	return responses.OK(ctx, resources.AuthToken{AccessToken: token, TokenType: "Bearer", ExpiresIn: tokenTTLSeconds()})
}

func (c *AuthController) Logout(ctx http.Context) http.Response {
	if err := facades.Auth(ctx).Logout(); err != nil {
		return responses.Error(ctx, err)
	}
	return responses.NoContent(ctx)
}

func (c *AuthController) issueToken(ctx http.Context, user *models.User) (resources.AuthToken, error) {
	token, err := facades.Auth(ctx).Login(user)
	if err != nil {
		return resources.AuthToken{}, err
	}

	me := c.presenter.Me(user)
	return resources.AuthToken{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   tokenTTLSeconds(),
		User:        &me,
	}, nil
}

func tokenTTLSeconds() int {
	return facades.Config().GetInt("jwt.ttl") * 60
}

func isExpired(err error) bool {
	return errors.Is(err, frameworkerrors.AuthTokenExpired)
}
