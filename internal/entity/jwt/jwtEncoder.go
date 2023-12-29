package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"time"
)

const (
	secret = "secretKey"
)

type JWT struct {
	tokenTTL time.Duration
	log      *slog.Logger
}

func NewJWT(tokenTTL time.Duration, log *slog.Logger) *JWT {
	return &JWT{
		tokenTTL: tokenTTL,
		log:      log,
	}
}

func (e JWT) NewToken(login string) (string, error) {
	const op = "JWT.NewToken"
	log := e.log.With(slog.String("op", op))

	log.Debug("starting generate token")
	token := jwt.New(jwt.SigningMethodHS256)

	log.Debug("init claims")
	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = login
	claims["exp"] = time.Now().Add(e.tokenTTL).Unix()

	log.Debug("claims initialized")
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Debug("cant generate token: " + err.Error())
		return "", err
	}

	return tokenString, nil
}
