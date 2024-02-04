package src

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
			r.PathPrefix(config.Get("passport.path", "oauth")).Namespace("goral/Passport/Http/Controllers").LoadRoutesFrom(currentDirectoryPath() + "/../routes/web.go")
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

func (s *PassportServiceProvider) RegisterCommands() {
	if app.IsRunningInConsole() {
		app.Commands(
			InstallCommand{},
			ClientCommand{},
			HashCommand{},
			KeysCommand{},
			PurgeCommand{},
		)
	}
}

func (s *PassportServiceProvider) Register() {
	// Merge the configuration from the `passport.php` file into the application's configuration
	app.MergeConfigFrom(currentDirectoryPath()+"/../config/passport.php", "passport")

	// Set whether to use client UUIDs based on the configuration
	Passport.SetClientUuids(app.Make(Config.Class).GetBool("passport.client_uuids"))

	// When the `AuthorizationController` is requested, inject the `StatefulGuard` with the configured guard
	app.When(AuthorizationController{}, func(context *application.Context) *auth.StatefulGuard {
		return auth.Guard(config.Get("passport.guard", nil))
	})

	// Register the OAuth authorization server
	s.RegisterAuthorizationServer()

	// Register the OAuth client repository
	s.RegisterClientRepository()

	// Register the JWT parser for authenticating OAuth tokens
	s.RegisterJWTParser()

	// Register the OAuth resource server for handling protected API requests
	s.RegisterResourceServer()

	// Register the OAuth guard for protecting routes and middleware
	s.RegisterGuard()

	// Set the authorization view for the OAuth authorization endpoint
	Passport.AuthorizationView("passport::authorize")
}

func (s *PassportServiceProvider) RegisterAuthorizationServer() {
	app.Singleton(AuthorizationServer.Class, func() *AuthorizationServer {
		server := s.MakeAuthorizationServer()

		server.DefaultScope = Passport.DefaultScope

		if err := server.Validate(); err != nil {
			log.Fatal("Error validating authorization server config:", err)
		}

		server.EnableGrantType(s.MakeAuthCodeGrant(), Passport.TokensExpireIn())
		server.EnableGrantType(s.MakeRefreshTokenGrant(), Passport.TokensExpireIn())
		server.EnableGrantType(s.MakePasswordGrant(), Passport.TokensExpireIn())
		server.EnableGrantType(s.PersonalAccessGrant(), Passport.PersonalAccessTokensExpireIn())
		server.EnableGrantType(s.ClientCredentialsGrant(), Passport.TokensExpireIn())
		if Passport.ImplicitGrantEnabled {
			server.EnableGrantType(s.MakeImplicitGrant(), Passport.TokensExpireIn())
		}

		return server
	})
}

func (s *PassportServiceProvider) makeAuthCodeGrant() (*AuthCodeGrant, error) {
	authCodeGrant, err := s.buildAuthCodeGrant()
	if err != nil {
		return nil, err
	}

	authCodeGrant.setRefreshTokenTTL(Passport.refreshTokensExpireIn())

	return authCodeGrant, nil
}

func (s *PassportServiceProvider) buildAuthCodeGrant() (*AuthCodeGrant, error) {
	authCodeRepository := app.Make(Bridge.AuthCodeRepository.Class)
	refreshTokenRepository := app.Make(Bridge.RefreshTokenRepository.Class)

	if authCodeRepository == nil || refreshTokenRepository == nil {
		return nil, errors.New("AuthCodeRepository or RefreshTokenRepository is not defined")
	}

	return newAuthCodeGrant(authCodeRepository, refreshTokenRepository, newDateInterval("PT10M")), nil
}

func (s *PassportServiceProvider) makeRefreshTokenGrant() (*RefreshTokenGrant, error) {
	refreshTokenGrant, err := s.buildRefreshTokenGrant()
	if err != nil {
		return nil, err
	}

	refreshTokenGrant.SetRefreshTokenTTL(Passport.refreshTokensExpireIn())

	return refreshTokenGrant, nil
}

func (s *PassportServiceProvider) buildRefreshTokenGrant() (*RefreshTokenGrant, error) {
	refreshTokenRepository := app.Make(RefreshTokenRepository.Class)

	if refreshTokenRepository == nil {
		return nil, errors.New("RefreshTokenRepository is not defined")
	}

	return newRefreshTokenGrant(refreshTokenRepository), nil
}

func (s *PassportServiceProvider) makePasswordGrant() (*PasswordGrant, error) {
	passwordGrant, err := s.buildPasswordGrant()
	if err != nil {
		return nil, err
	}

	passwordGrant.SetRefreshTokenTTL(Passport.refreshTokensExpireIn())

	return passwordGrant, nil
}

func (s *PassportServiceProvider) buildPasswordGrant() (*PasswordGrant, error) {
	userRepository := app.Make(Bridge.UserRepository.Class)
	refreshTokenRepository := app.Make(Bridge.RefreshTokenRepository.Class)

	if userRepository == nil || refreshTokenRepository == nil {
		return nil, errors.New("UserRepository or RefreshTokenRepository is not defined")
	}

	return newPasswordGrant(userRepository, refreshTokenRepository), nil
}

func (s *PassportServiceProvider) makeImplicitGrant() (*ImplicitGrant, error) {
	implicitGrant, err := newImplicitGrant(Passport.tokensExpireIn())
	if err != nil {
		return nil, err
	}

	return implicitGrant, nil
}

func (s *PassportServiceProvider) makeAuthorizationServer() (*AuthorizationServer, error) {
	clientRepository := app.Make(Bridge.ClientRepository.Class)
	accessTokenRepository := app.Make(Bridge.AccessTokenRepository.Class)
	scopeRepository := app.Make(Bridge.ScopeRepository.Class)
	encryptionKey := s.makeCryptKey("private")

	if clientRepository == nil || accessTokenRepository == nil || scopeRepository == nil {
		return nil, errors.New("ClientRepository, AccessTokenRepository, or ScopeRepository is not defined")
	}

	if encryptionKey == nil {
		return nil, errors.New("Encryption key is not defined")
	}

	return newAuthorizationServer(clientRepository, accessTokenRepository, scopeRepository, encryptionKey, app.Make("encrypter").GetKey(), Passport.authorizationServerResponseType), nil
}

func (s *PassportServiceProvider) registerClientRepository() {
	config := app.Make("config").Get("passport.personal_access_client")
	clientId := config["id"]
	clientSecret := config["secret"]

	if clientId == "" || clientSecret == "" {
		err := errors.New("Personal access client ID or secret is not defined")
		logs.Print("Error:", err)
		return
	}

	app.Singleton(ClientRepository.Class, func() *ClientRepository {
		return newClientRepository(clientId, clientSecret)
	})
}

func (s *PassportServiceProvider) registerJWTParser() {
	if err := app.Register(ParserContract.Class, func() Parser {
		return newParser(newJoseEncoder)
	}); err != nil {
		logs.Print("Error:", err)
		return
	}
}


func (s *PassportServiceProvider) registerResourceServer() {
    accessTokenRepository := app.Make(Bridge.AccessTokenRepository.Class)
    encryptionKey := s.makeCryptKey("public")

    if accessTokenRepository == nil || encryptionKey == nil {
        err := errors.New("AccessTokenRepository or encryption key is not defined")
        logs.Print("Error:", err)
        return
    }

    app.Singleton(ResourceServer.Class, func () *ResourceServer {
        return newResourceServer(accessTokenRepository, encryptionKey)
    })
}

func (s *PassportServiceProvider) makeCryptKey(datatype string) string {
	key := str_replace("\\n", "\n", app.Make("config").Get("passport."+datatype+"_key"))

	if key == "" {
		key = "file://" + Passport.keyPath("oauth-"+datatype+".key")
	}

	return newCryptKey(key, nil, false)
}


func (s *PassportServiceProvider) registerGuard() {
    Auth.Resolved(func (auth *Auth) {
        auth.Extend("passport", func (app *application, name string, config map[string]interface{}) (*Guard, error) {
            guard, err := s.makeGuard(config)
            if err != nil {
                return nil, err
            }

            app.Refresh("request", guard, "setRequest")
            return guard, nil
        })
    })
}

func (s *PassportServiceProvider) makeGuard(config map[string]interface{}) (*Guard, error) {
    resourceServer := app.Make(ResourceServer.Class)
    userProvider := newPassportUserProvider(Auth.CreateUserProvider(config["provider"]), config["provider"])
    tokenRepository := app.Make(TokenRepository.Class)
    clientRepository := app.Make(ClientRepository.Class)
    encrypter := app.Make("encrypter")
    request := app.Make("request")

    if resourceServer == nil || userProvider == nil || tokenRepository == nil || clientRepository == nil || encrypter == nil || request == nil {
        return nil, errors.New("ResourceServer, UserProvider, TokenRepository, ClientRepository, Encrypter, or Request is not defined")
    }

    return newGuard(resourceServer, userProvider, tokenRepository, clientRepository, encrypter, request), nil
}

func (s *PassportServiceProvider) deleteCookieOnLogout() {
    Event.Listen(Logout.Class, func () {
        if request.HasCookie(Passport.cookie()) {
            Cookie.Queue(Cookie.Forget(Passport.cookie()))
        }
    })
}
