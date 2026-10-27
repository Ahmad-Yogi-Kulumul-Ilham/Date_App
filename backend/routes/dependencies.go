package routes

import (
	"time"

	"leochat/app/facades"
	"leochat/app/http/controllers"
	"leochat/app/http/resources"
	"leochat/app/repositories"
	"leochat/app/services"
)

type apiControllers struct {
	auth      *controllers.AuthController
	profile   *controllers.ProfileController
	discovery *controllers.DiscoveryController
	swipe     *controllers.SwipeController
	match     *controllers.MatchController
}

// newAPIControllers wires repositories into services into controllers. This is
// the only place that knows about concrete implementations.
func newAPIControllers() apiControllers {
	now := services.Clock(func() time.Time { return time.Now().UTC() })
	orm := facades.Orm()
	publicDisk := facades.Storage().Disk("public")

	users := repositories.NewUserRepository(orm)
	photos := repositories.NewPhotoRepository(orm)
	candidates := repositories.NewCandidateRepository(orm)
	swipes := repositories.NewSwipeRepository(orm)
	matches := repositories.NewMatchRepository(orm)
	messages := repositories.NewMessageRepository(orm)

	presenter := resources.NewPresenter(publicDisk.Url, now)

	return apiControllers{
		auth:      controllers.NewAuthController(services.NewAuthService(users, facades.Hash(), now), presenter),
		profile:   controllers.NewProfileController(services.NewProfileService(users, photos, publicDisk, now), presenter),
		discovery: controllers.NewDiscoveryController(services.NewDiscoveryService(users, candidates, now), presenter),
		swipe:     controllers.NewSwipeController(services.NewSwipeService(users, swipes, matches)),
		match:     controllers.NewMatchController(services.NewChatService(matches, messages, now), presenter),
	}
}
