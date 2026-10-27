package middleware

import (
	"errors"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	frameworkerrors "github.com/goravel/framework/errors"

	"leochat/app/facades"
	"leochat/app/http/responses"
)

type contextKey string

const userIDKey contextKey = "auth.user_id"

// Authenticate rejects requests without a valid bearer token and stores the
// authenticated user's ID on the context for CurrentUserID.
type Authenticate struct{}

func (m *Authenticate) Signature() string { return "leochat:authenticate" }

func (m *Authenticate) Handle(ctx http.Context) {
	token := ctx.Request().Header("Authorization")
	if token == "" {
		_ = responses.Fail(ctx, http.StatusUnauthorized, "unauthenticated", "Silakan login terlebih dahulu.").Abort()
		return
	}

	payload, err := facades.Auth(ctx).Parse(token)
	if errors.Is(err, frameworkerrors.AuthTokenExpired) {
		_ = responses.Fail(ctx, http.StatusUnauthorized, "token_expired", "Sesi habis, silakan perbarui token.").Abort()
		return
	}
	if err != nil {
		_ = responses.Fail(ctx, http.StatusUnauthorized, "unauthenticated", "Token tidak valid.").Abort()
		return
	}

	id, err := strconv.ParseUint(payload.Key, 10, 64)
	if err != nil {
		_ = responses.Fail(ctx, http.StatusUnauthorized, "unauthenticated", "Token tidak valid.").Abort()
		return
	}

	ctx.WithValue(userIDKey, uint(id))
	ctx.Request().Next()
}

// CurrentUserID returns the ID set by Authenticate. Only call it on routes
// behind that middleware.
func CurrentUserID(ctx http.Context) uint {
	id, _ := ctx.Value(userIDKey).(uint)
	return id
}
