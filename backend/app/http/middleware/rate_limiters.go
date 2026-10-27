package middleware

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/http/limit"

	"leochat/app/facades"
	"leochat/app/http/responses"
)

const (
	LimiterAuth  = "auth"
	LimiterSwipe = "swipe"
	LimiterChat  = "chat"
)

// RegisterRateLimiters defines the named limiters used by the throttle middleware.
func RegisterRateLimiters() {
	limiter := facades.RateLimiter()

	limiter.For(LimiterAuth, func(ctx http.Context) http.Limit {
		return limit.PerMinute(10).By(ctx.Request().Ip()).Response(tooManyRequests)
	})

	limiter.For(LimiterSwipe, func(ctx http.Context) http.Limit {
		return limit.PerMinute(120).By(userKey(ctx)).Response(tooManyRequests)
	})

	limiter.For(LimiterChat, func(ctx http.Context) http.Limit {
		return limit.PerMinute(60).By(userKey(ctx)).Response(tooManyRequests)
	})
}

func userKey(ctx http.Context) string {
	return "user:" + strconv.FormatUint(uint64(CurrentUserID(ctx)), 10)
}

func tooManyRequests(ctx http.Context) {
	_ = responses.Fail(ctx, http.StatusTooManyRequests, "too_many_requests", "Terlalu banyak permintaan. Tunggu sebentar.").Abort()
}
