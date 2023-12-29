package encryptorImpl

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type BcryptEncryptor struct {
	log *slog.Logger
}

func NewBcryptEncryptor(log *slog.Logger) *BcryptEncryptor {
	return &BcryptEncryptor{
		log: log,
	}
}

func (e BcryptEncryptor) Encrypt(password string) (string, error) {
	const op = "BcryptEncryptor.Encrypt"
	log := e.log.With(slog.String("op", op))

	log.Debug("start encrypt")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Debug("cant encrypt: " + err.Error())
		return "", err
	}

	return string(hashedPassword), nil
}

func (e BcryptEncryptor) Compare(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Println("password: " + password + " hashed: " + hashedPassword)
		fmt.Println(err.Error())
		return false
	}

	return true
}
