package routes

import (
	"github.com/goravel/framework/contracts/http"

	"leochat/app/facades"
)

func Web() {
	facades.Route().Get("/health", func(ctx http.Context) http.Response {
		return ctx.Response().Success().Json(http.Json{"status": "ok"})
	})

	// Uploaded photos on the "public" disk. In production serve these from object storage or a CDN.
	facades.Route().Static("/storage", "./storage/app/public")
}
