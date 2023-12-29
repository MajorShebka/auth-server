package entity

type Customer struct {
	Login    string `bson:"login"`
	Password string `bson:"password"`
}
