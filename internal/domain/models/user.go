package models

type User struct {
	Id    string `bson:"-"`
	Email string `bson:"email"`
	UUID  string `bson:"uuid"`
}