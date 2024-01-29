package src

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ApiTokenCookieFactory struct {
	config    *Config
	encrypter Encrypter
}

func NewApiTokenCookieFactory(config *Config, encrypter Encrypter) *ApiTokenCookieFactory {
	return &ApiTokenCookieFactory{
		config:    config,
		encrypter: encrypter,
	}
}

func (f *ApiTokenCookieFactory) Make(userId string, csrfToken string) *http.Cookie {
	config := f.config.Get("session")

	expiration := time.Now().Add(time.Duration(config.Lifetime) * time.Minute)

	return &http.Cookie{
		Name:     Passport.Cookie(),
		Value:    f.CreateToken(userId, csrfToken, expiration),
		Expires:  expiration,
		Path:     config.Path,
		Domain:   config.Domain,
		Secure:   config.Secure,
		HttpOnly: true,
		SameSite: config.SameSite,
	}
}

func (f *ApiTokenCookieFactory) CreateToken(userId string, csrfToken string, expiration time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    userId,
		"csrf":   csrfToken,
		"expiry": expiration.Unix(),
	})

	key, err := f.encrypter.PublicKey()
	if err != nil {
		panic(err)
	}

	return token.SignedString(key)
}
