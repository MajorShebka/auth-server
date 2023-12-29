package impl

import (
	"authServer/internal/DTO"
	"authServer/internal/entity"
	"authServer/internal/errors/repositoryErrors"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"time"
)

const (
	databaseName   = "authServer"
	collectionName = "customers"
)

type MongoCustomerRepository struct {
	url string
	log *slog.Logger
}

func NewMongoCustomerRepo(url string, log *slog.Logger) *MongoCustomerRepository {
	return &MongoCustomerRepository{
		url: url,
		log: log,
	}
}

func (repo MongoCustomerRepository) getMongoCollection() (*mongo.Collection, context.CancelFunc) {
	const op = "MongoCustomerRepository.getMongoCollection"
	log := repo.log.With(
		slog.String("op", op),
	)

	log.Debug("opening connection")
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(repo.url))
	if err != nil {
		panic(err)
	}

	db := client.Database(databaseName)
	collection := db.Collection(collectionName)

	log.Debug("ending op")

	return collection, cancel
}

func (repo MongoCustomerRepository) FindCustomerByName(customerLogin string) (entity.Customer, error) {
	const op = "MongoCustomerRepository.FindCustomerByName"
	log := repo.log.With(
		slog.String("op", op),
	)
	collection, cancel := repo.getMongoCollection()
	defer cancel()

	log.Debug("starting finding customer")
	cursor := collection.FindOne(context.TODO(), bson.M{"login": customerLogin})

	var foundCustomer entity.Customer

	err := cursor.Decode(&foundCustomer)
	log.Debug("found customer: ", foundCustomer)
	if err != nil {
		log.Debug("ending op with error " + err.Error())
		return entity.Customer{}, repositoryErrors.CustomerNotFoundErr{}
	}

	repo.log.Debug(op + " ending op")
	return foundCustomer, nil
}

func (repo MongoCustomerRepository) FindCustomer(customer DTO.CustomerDTO) (entity.Customer, error) {
	const op = "MongoCustomerRepository.FindCustomer"
	collection, cancel := repo.getMongoCollection()
	defer cancel()

	repo.log.Debug(op + " starting finding customer")
	cursor := collection.FindOne(context.TODO(), customer)

	var foundCustomer entity.Customer

	err := cursor.Decode(&foundCustomer)
	repo.log.Debug(op+" found customer: ", foundCustomer)
	if err != nil {
		repo.log.Debug(op + " ending op with error " + err.Error())
		return entity.Customer{}, repositoryErrors.CustomerNotFoundErr{}
	}

	repo.log.Debug(op + " ending op")

	return foundCustomer, nil
}

func (repo MongoCustomerRepository) CreateCustomer(customer entity.Customer) error {
	collection, cancel := repo.getMongoCollection()
	defer cancel()

	_, err := collection.InsertOne(context.TODO(), customer)
	if err != nil {
		repo.log.Debug("cant insert: " + err.Error())
		return errors.New("customer already exists")
	}

	return nil
}
