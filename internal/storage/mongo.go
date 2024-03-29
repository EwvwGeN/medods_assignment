package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/EwvwGeN/medods_assignment/internal/config"
	"github.com/EwvwGeN/medods_assignment/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoProvider struct {
	cfg config.MongoConfig
	db  *mongo.Database
}

func NewMongoProvider(ctx context.Context, cfg config.MongoConfig) (*mongoProvider, error) {
	client, err := mongo.Connect(ctx,
		options.Client().
			ApplyURI(
				fmt.Sprintf("%s://%s:%s@%s:%s/?authSource=%s",
					cfg.ConectionFormat,
					cfg.User,
					cfg.Password,
					cfg.Host,
					cfg.Port,
					cfg.AuthSourse,
				)),
	)
	if err != nil {
		return nil, err
	}
	dbList, err := client.ListDatabases(ctx, bson.D{{Key: "name", Value: cfg.Database}})
	if err != nil {
		return nil, err
	}
	if dbList.TotalSize == 0 {
		return nil, ErrDbNotExist
	}
	db := client.Database(cfg.Database)
	colList, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: cfg.UserCollection}})
	if err != nil {
		return nil, err
	}
	if len(colList) != 1 {
		return nil, ErrCollNotExist
	}
	return &mongoProvider{
		cfg: cfg,
		db:  db,
	}, nil
}

func (m *mongoProvider) SaveUser(ctx context.Context, email, uuid string) (err error) {
	_, err = m.db.Collection(m.cfg.UserCollection).InsertOne(ctx, bson.D{
		{Key: "email", Value: email},
		{Key: "uuid", Value: uuid},
	})
	if mongo.IsDuplicateKeyError(err) {
		return ErrUserExist
	}
	return
}

func (m *mongoProvider) GetUserByUUID(ctx context.Context, uuid string) (user *models.User, err error) {
	findedUser := m.db.Collection(m.cfg.UserCollection).FindOne(ctx, bson.D{{Key: "uuid", Value: uuid}})
	if raw, _ := findedUser.Raw(); raw != nil {
		findedUser.Decode(&user)
		return user, nil
	}
	return &models.User{}, ErrUserNotFound
}

func (m *mongoProvider) SaveRefreshToken(ctx context.Context, uuid, refresh string, refreshTTL time.Duration) (err error) {
	findedUser := m.db.Collection(m.cfg.UserCollection).FindOne(ctx, bson.D{{Key: "uuid", Value: uuid}})
	if raw, _ := findedUser.Raw(); raw == nil {
		return ErrUserNotFound
	}
	_, err = m.db.Collection(m.cfg.UserCollection).UpdateOne(ctx, bson.D{
		{Key: "uuid", Value: uuid},
	},
	bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "refresh_token", Value: refresh},
			{Key: "expires_at", Value: time.Now().Add(refreshTTL).Unix()}},
		},
	})
	if err != nil {
		return ErrUpdate
	}
	return
}
