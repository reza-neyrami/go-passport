package main

import (

	console "github.com/reza-neyrami/go-passport/src/Console"
	guards "github.com/reza-neyrami/go-passport/src/Guards"
)

type PassportServiceProvider struct {
}

func (s *PassportServiceProvider) Boot(app *application.Application) {
	// Register the routes
	if config.Get("passport.registersRoutes") {
		app.router.group(config.Get("passport.path", "oauth").As("passport"), func(r *router.Router) {
			r.get("/authorize", controllers.AuthorizeController.Handle())
			r.post("/authorize", controllers.AuthorizeController.Handle())
			r.get("/token", controllers.AccessTokenController.Handle())
			r.post("/token", controllers.AccessTokenController.Handle())
			r.get("/userinfo", controllers.UserInfoController.Handle())
		})
	}

	// Register the resources
	app.views.AddPath("passport", config.Get("passport.resourcesPath", "resources/views/passport"))

	// Register the migration files
	if app.runningInConsole() && config.Get("passport.runsMigrations") && !config.Get("passport.clientUuids") {
		app.migrations.AddPath(config.Get("passport.migrationsPath", "database/migrations"))
	}

	// Register the publishable resources
	if app.runningInConsole() {
		app.publish(config.Get("passport.migrationsPath", "database/migrations"), "passport-migrations")
		app.publish(config.Get("passport.configPath", "config/passport.php"), "passport-config")
		app.publish(config.Get("passport.resourcesPath", "resources/views/passport"), "passport-views")
	}

	// Register the commands
	app.commands.Add(
		console.InstallCommand,
		console.ClientCommand,
		console.HashCommand,
		console.KeysCommand,
		console.PurgeCommand,
	)

	// Register the guard
	app.auth.Extend("passport", func(auth *auth.Guard) *guards.TokenGuard {
		return guards.NewTokenGuard(
			auth,
			models.NewClientRepository(),
			models.NewAccessTokenRepository(),
			models.NewRefreshTokenRepository(),
			models.NewScopeRepository(),
			jwt.Parser{}, // Added import and initialized jwt.Parser
			app.make("encrypter"),
			app.make("request"),
		)
	})

	// Delete the cookie on logout
	if config.Get("passport.deleteCookieOnLogout") {
		app.Events.Listen(events.Logout, func(event *events.LogoutEvent) {
			if event.User != nil {
				if cookie := event.User.GetToken(); cookie != nil {
					cookie.Delete()
				}
			}
		})
	}
}

func (s *PassportServiceProvider) RegisterRoutes() {
	if Passport.RegistersRoutes {
		router := app.NewRouter()
		router.Group("passport", func(r *router.Router) {
			r.PathPrefix(config.Get("passport.path", "oauth")).Namespace("Laravel/Passport/Http/Controllers").LoadRoutesFrom(currentDirectoryPath() + "/../routes/web.go")
		})
		app.Router.Merge(router)
	}
}

func (s *PassportServiceProvider) RegisterResources() {
	app.Views.AddPath("passport", config.Get("passport.resourcesPath", "resources/views/passport"))
}

func (s *PassportServiceProvider) RegisterMigrations() {
	if app.IsRunningInConsole() && Passport.RunsMigrations && !config.GetBool("passport.clientUuids") {
		app.Migrations.AddPath(config.Get("passport.migrationsPath", "database/migrations"))
	}
}
