package itemRepositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ItemRepositoriesService interface{}

	itemRepository struct {
		db *mongo.Client
	}
)
func NewItemRepository (db *mongo.Client) ItemRepositoriesService{
	return &itemRepository{
		db: db,
	}
}
func (r *itemRepository) itemDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("item_db")
}