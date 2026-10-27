package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"leochat/app/http/middleware"
	"leochat/app/http/requests"
	"leochat/app/http/resources"
	"leochat/app/http/responses"
	"leochat/app/services"
)

type MatchController struct {
	chat      *services.ChatService
	presenter *resources.Presenter
}

func NewMatchController(chat *services.ChatService, presenter *resources.Presenter) *MatchController {
	return &MatchController{chat: chat, presenter: presenter}
}

func (c *MatchController) Index(ctx http.Context) http.Response {
	userID := middleware.CurrentUserID(ctx)

	summaries, err := c.chat.ListMatches(ctx, userID)
	if err != nil {
		return responses.Error(ctx, err)
	}
	return responses.OK(ctx, c.presenter.MatchSummaries(summaries, userID))
}

// Messages pages backwards through a conversation: pass the smallest id you
// have as before_id to load older messages.
func (c *MatchController) Messages(ctx http.Context) http.Response {
	matchID, ok := routeID(ctx, "id")
	if !ok {
		return notFound(ctx)
	}

	userID := middleware.CurrentUserID(ctx)
	beforeID := ctx.Request().QueryInt64("before_id")
	if beforeID < 0 {
		beforeID = 0
	}

	messages, err := c.chat.Messages(ctx, userID, matchID, uint(beforeID), ctx.Request().QueryInt("limit"))
	if err != nil {
		return responses.Error(ctx, err)
	}
	return responses.OK(ctx, c.presenter.Messages(messages, userID))
}

func (c *MatchController) SendMessage(ctx http.Context) http.Response {
	matchID, ok := routeID(ctx, "id")
	if !ok {
		return notFound(ctx)
	}

	var req requests.SendMessageRequest
	if res, ok := validate(ctx, &req); !ok {
		return res
	}

	userID := middleware.CurrentUserID(ctx)
	message, err := c.chat.Send(ctx, userID, matchID, req.Body)
	if err != nil {
		return responses.Error(ctx, err)
	}
	return responses.Created(ctx, c.presenter.Message(message, userID))
}
