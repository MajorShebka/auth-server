package encryptor

import (
	"authServer/internal/entity/encryptor/encryptorImpl"
	"log/slog"
)

type Encryptor interface {
	Encrypt(password string) (string, error)
	Compare(hashedPassword string, password string) bool
}

func InitEncryptor(log *slog.Logger) *Encryptor {
	var encryptor Encryptor
	encryptor = encryptorImpl.NewBcryptEncryptor(log)
	return &encryptor
}
