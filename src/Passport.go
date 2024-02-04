package src

import (
	"crypto/cipher"
	"time"
)

type Passport struct {
	ImplicitGrantEnabled         bool
	DefaultScope                 string
	Scopes                       []Scope
	TokensExpireIn               time.Duration
	RefreshTokensExpireIn        time.Duration
	PersonalAccessTokensExpireIn time.Duration
	Cookie                       string

	IgnoreCsrfToken                 bool
	KeyPath                         string
	AccessTokenEntity               string
	AuthCodeModel                   string
	ClientModel                     string
	ClientUuids                     bool
	PersonalAccessClientModel       string
	TokenModel                      string
	RefreshTokenModel               string
	HashesClientSecrets             bool
	TokenEncryptionKeyCallback      func(encrypter *cipher.AEAD) []byte
	AuthorizationView               func() string
	WithInheritedScopes             bool
	AuthorizationServerResponseType string
	RegistersRoutes                 bool
}



func NewPassport() *Passport {
	return &Passport{
		ImplicitGrantEnabled:         false,
		DefaultScope:                 "api",
		Scopes:                       []Scope{},
		TokensExpireIn:               time.Hour * 24,
		RefreshTokensExpireIn:        time.Hour * 3600,
		PersonalAccessTokensExpireIn: time.Hour * 3600,
		Cookie:                       "laravel_token",

		IgnoreCsrfToken:           false,
		KeyPath:                   "storage/oauth-keys",
		AccessTokenEntity:         "App/Models/AccessToken",
		AuthCodeModel:             "App/Models/AuthCode",
		ClientModel:               "App/Models/Client",
		ClientUuids:               true,
		PersonalAccessClientModel: "App/Models/PersonalAccessClient",
		TokenModel:                "App/Models/Token",
		RefreshTokenModel:         "App/Models/RefreshToken",
		HashesClientSecrets:       false,
		TokenEncryptionKeyCallback: func(encrypter *cipher.AEAD) []byte {
			return nil
		},
		AuthorizationView: func() string {
			return "/oauth/authorize"
		},
		WithInheritedScopes:             false,
		AuthorizationServerResponseType: "html",
		RegistersRoutes:                 true,
	}
}

func (p *Passport) EnableImplicitGrant() *Passport {
	p.ImplicitGrantEnabled = true
	return p
}

func (p *Passport) setDefaultScope(scope string) {
	p.DefaultScope = scope
}

func (p *Passport) scopeIds() []string {
	return []string{}
}

func (p *Passport) hasScope(id string) bool {
	return false
}

func (p *Passport) scopes() []Scope {
	return nil
}

func (p *Passport) scopesFor(ids []string) []Scope {
	return nil
}

func (p *Passport) tokensCan(scopes []Scope) *Passport {
	p.Scopes = scopes
	return p
}

func (p *Passport) tokensExpireIn(date time.Time) *Passport {
	p.TokensExpireIn = time.Duration(date.Sub(time.Now()))
	return p
}

func (p *Passport) refreshTokensExpireIn(date time.Time) *Passport {
	p.RefreshTokensExpireIn = time.Duration(date.Sub(time.Now()))
	return p
}

func (p *Passport) personalAccessTokensExpireIn(date time.Time) *Passport {
	p.PersonalAccessTokensExpireIn = time.Duration(date.Sub(time.Now()))
	return p
}

func (p *Passport) cookie(cookie string) string {
	if cookie == "" {
		return p.Cookie
	}
	return cookie
}

func (p *Passport) ignoreCsrfToken(ignore bool) *Passport {
	p.IgnoreCsrfToken = ignore
	return p
}

func (p *Passport) useAccessTokenEntity(entity string) *Passport {
	p.AccessTokenEntity = entity
	return p
}

func (p *Passport) useAuthCodeModel(model string) *Passport {
	p.AuthCodeModel = model
	return p
}

func (p *Passport) useClientModel(model string) *Passport {
	p.ClientModel = model
	return p
}

func (p *Passport) usePersonalAccessClientModel(model string) *Passport {
	p.PersonalAccessClientModel = model
	return p
}

func (p *Passport) useTokenModel(model string) *Passport {
	p.TokenModel = model
	return p
}

func (p *Passport) useRefreshTokenModel(model string) *Passport {
	p.RefreshTokenModel = model
	return p
}

func (p *Passport) hashClientSecrets(hash bool) *Passport {
	p.HashesClientSecrets = hash
	return p
}

func (p *Passport) encryptTokensUsing(callback func(encrypter *cipher.AEAD) []byte) *Passport {
	p.TokenEncryptionKeyCallback = callback
	return p
}

func (p *Passport) authorizationView(view func() string) *Passport {
	p.AuthorizationView = view
	return p
}

func (p *Passport) ignoreRoutes() *Passport {
	p.RegistersRoutes = false
	return p
}

func (p *Passport) ignoreMigrations() *Passport {
	p.RunsMigrations = false
	return p
}

func (p *Passport) withCookieSerialization() *Passport {
	p.UnserializesCookies = true
	return p
}

func (p *Passport) withoutCookieSerialization() *Passport {
	p.UnserializesCookies = false
	return p
}

func (p *Passport) withCookieEncryption() *Passport {
	p.DecryptsCookies = true
	return p
}

func (p *Passport) withoutCookieEncryption() *Passport {
	p.DecryptsCookies = false
	return p
}

func (p *Passport) setTokenEncryptionKey(key []byte) *Passport {
	p.TokenEncryptionKey = key
	return p
}

func (p *Passport) getNonce() []byte {
	return p.TokenEncryptionKeyCallback(p.TokenEncryptionKey)
}
