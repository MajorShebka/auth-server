package DTO

type CustomerDTO struct {
	Login    string `bson:"login"`
	Password string `bson:"password"`
}
