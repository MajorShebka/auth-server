package auth

import (
	"authServer/internal/DTO"
	"authServer/internal/entity/encryptor"
	"authServer/internal/entity/jwt"
	"authServer/internal/repository/customerRepository"
	"authServer/internal/service/auth/impl"
	"log/slog"
	"time"
)

type Auth interface {
	Login(customerDTO DTO.CustomerDTO) (string, error)
	Register(customerDTO DTO.CustomerDTO) (string, error)
}

func InitAuth(repo *customerRepository.CustomerRepository, log *slog.Logger, tokenTTL time.Duration) *Auth {
	jwt := jwt.NewJWT(tokenTTL, log)
	enc := encryptor.InitEncryptor(log)

	var auth Auth
	auth = impl.NewAuthImpl(repo, log, jwt, enc)

	return &auth
}
