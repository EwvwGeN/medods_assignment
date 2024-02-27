package models

type User struct {
	Id          string `bson:"-"`
	Email       string `bson:"email"`
	UUID        string `bson:"uuid"`
	RefreshHash string `bson:"refresh_token"`
	ExpiresAt   int64  `bson:"expires_at"`
}