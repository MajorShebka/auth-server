package customerRepository

import (
	"authServer/internal/DTO"
	"authServer/internal/entity"
	"authServer/internal/repository/customerRepository/impl"
	"log/slog"
)

type CustomerRepository interface {
	FindCustomerByName(customerName string) (entity.Customer, error)
	FindCustomer(customerDTO DTO.CustomerDTO) (entity.Customer, error)
	CreateCustomer(customer entity.Customer) error
}

func InitCustomerRepository(storageUrl string, log *slog.Logger) *CustomerRepository {
	var customerRepository CustomerRepository
	customerRepository = impl.NewMongoCustomerRepo(storageUrl, log)

	return &customerRepository
}
