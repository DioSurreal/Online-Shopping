package inventoryRepositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type(
	InventoryRepositoriesService interface{}

	inventoryRepository struct{
		db  *mongo.Client
	}
)

func NewInventoryRepository (db *mongo.Client) InventoryRepositoriesService {
	return inventoryRepository{db}
}

func (r *inventoryRepository) inventoryDbConn(pctx *context.Context) *mongo.Database {
	return r.db.Database("inventory")
}