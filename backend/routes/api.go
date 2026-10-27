package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/http/middleware"

	"leochat/app/facades"
	appmiddleware "leochat/app/http/middleware"
)

func Api() {
	c := newAPIControllers()
	authenticate := &appmiddleware.Authenticate{}

	facades.Route().Prefix("/api/v1").Group(func(r route.Router) {
		r.Prefix("/auth").Middleware(middleware.Throttle(appmiddleware.LimiterAuth)).Group(func(r route.Router) {
			r.Post("/register", c.auth.Register)
			r.Post("/login", c.auth.Login)
			r.Post("/refresh", c.auth.Refresh)
		})

		r.Middleware(authenticate).Group(func(r route.Router) {
			r.Post("/auth/logout", c.auth.Logout)

			r.Get("/me", c.profile.Show)
			r.Patch("/me", c.profile.Update)
			r.Put("/me/location", c.profile.UpdateLocation)
			r.Post("/me/photos", c.profile.UploadPhoto)
			r.Delete("/me/photos/{id}", c.profile.DeletePhoto)

			r.Get("/discover", c.discovery.Index)
			r.Middleware(middleware.Throttle(appmiddleware.LimiterSwipe)).Post("/swipes", c.swipe.Store)

			r.Get("/matches", c.match.Index)
			r.Get("/matches/{id}/messages", c.match.Messages)
			r.Middleware(middleware.Throttle(appmiddleware.LimiterChat)).Post("/matches/{id}/messages", c.match.SendMessage)
		})
	})
}
