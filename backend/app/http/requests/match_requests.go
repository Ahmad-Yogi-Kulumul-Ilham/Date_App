package requests

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"leochat/app/models"
)

type SwipeRequest struct {
	UserID    int    `form:"user_id" json:"user_id"`
	Direction string `form:"direction" json:"direction"`
}

func (r *SwipeRequest) Authorize(http.Context) error { return nil }

func (r *SwipeRequest) Rules(http.Context) map[string]any {
	return map[string]any{
		"user_id":   "required|integer|min:1",
		"direction": "required|in:like,pass",
	}
}

func (r *SwipeRequest) Target() uint { return uint(r.UserID) }

func (r *SwipeRequest) SwipeDirection() models.SwipeDirection {
	return models.SwipeDirection(r.Direction)
}

const maxMessageLength = 2000

type SendMessageRequest struct {
	Body string `form:"body" json:"body"`
}

func (r *SendMessageRequest) Authorize(http.Context) error { return nil }

func (r *SendMessageRequest) Rules(http.Context) map[string]any {
	return map[string]any{
		"body": "required|string|max:" + itoa(maxMessageLength),
	}
}

func itoa(n int) string { return strconv.Itoa(n) }
