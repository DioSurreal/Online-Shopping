package userRepositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	UserRepositoriesService interface{}

	userRepository struct{
		db *mongo.Client
	}
)

func NewUserRepository(db *mongo.Client) UserRepositoriesService {
	return &userRepository{db: db}
}

func (r *userRepository) userDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("user")
}