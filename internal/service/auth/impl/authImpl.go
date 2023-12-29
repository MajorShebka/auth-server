package impl

import (
	"authServer/internal/DTO"
	"authServer/internal/entity"
	"authServer/internal/entity/encryptor"
	"authServer/internal/entity/jwt"
	"authServer/internal/errors/repositoryErrors"
	"authServer/internal/repository/customerRepository"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type AuthImpl struct {
	repo       *customerRepository.CustomerRepository
	log        *slog.Logger
	jwtEncoder *jwt.JWT
	enc        *encryptor.Encryptor
}

func NewAuthImpl(
	repo *customerRepository.CustomerRepository, log *slog.Logger, jwtEncoder *jwt.JWT, enc *encryptor.Encryptor,
) *AuthImpl {
	return &AuthImpl{
		repo:       repo,
		log:        log,
		jwtEncoder: jwtEncoder,
		enc:        enc,
	}
}

func (s AuthImpl) Login(customerDTO DTO.CustomerDTO) (string, error) {
	const op = "AuthImpl.Login()"
	log := s.log.With(slog.String("op", op))

	log.Debug("customer from params: ", customerDTO)
	customer, err := (*s.repo).FindCustomerByName(customerDTO.Login)
	if errors.Is(err, repositoryErrors.CustomerNotFoundErr{}) {
		log.Debug("ending with error: " + err.Error())
		return "", status.Error(codes.NotFound, "incorrect password or login")
	}
	log.Debug("found customer: ", customer)

	if !(*s.enc).Compare(customer.Password, customerDTO.Password) {
		log.Debug("incorrect password")
		return "", status.Error(codes.NotFound, "incorrect password or login")
	}

	token, err := s.createToken(customer.Login)
	if err != nil {
		log.Debug("Cant create token: " + err.Error())
		return "", status.Error(codes.Internal, "internal error")
	}

	return token, nil
}

func (s AuthImpl) Register(customerDTO DTO.CustomerDTO) (string, error) {
	const op = "AuthImpl.Login()"
	log := s.log.With(slog.String("op", op))

	log.Debug("customer from params: ", customerDTO)

	_, err := (*s.repo).FindCustomerByName(customerDTO.Login)
	if err == nil {
		log.Debug("customer already exist")
		return "", status.Error(codes.AlreadyExists, "already exists")
	}

	customer, err := s.initCustomer(customerDTO.Login, customerDTO.Password)
	if err != nil {
		log.Debug("cant init customer: " + err.Error())
		return "", status.Error(codes.Internal, "internal error")
	}

	err = (*s.repo).CreateCustomer(customer)
	if err != nil {
		log.Debug("error: " + err.Error())
		return "", status.Error(codes.Internal, "internal error")
	}

	token, err := s.createToken(customer.Login)
	if err != nil {
		log.Debug("Cant create token: " + err.Error())
		return "", status.Error(codes.Internal, "internal error")
	}

	return token, nil
}

func (s AuthImpl) initCustomer(customerLogin string, customerPassword string) (entity.Customer, error) {
	const op = "AuthImpl.initCustomer"
	log := s.log.With(slog.String("op", op))

	hashedPassword, err := (*s.enc).Encrypt(customerPassword)
	if err != nil {
		log.Debug("cant encrypt: " + err.Error())
		return entity.Customer{}, err
	}

	customer := entity.Customer{
		Login:    customerLogin,
		Password: hashedPassword,
	}

	return customer, nil
}

func (s AuthImpl) createToken(login string) (string, error) {
	const op = "authImpl.createToken"
	log := s.log.With(slog.String("op", op))
	token, err := (s.jwtEncoder).NewToken(login)
	if err != nil {
		log.Debug("cannot create token: " + err.Error())
		return "", status.Error(codes.Internal, "internal error")
	}

	return token, nil
}
