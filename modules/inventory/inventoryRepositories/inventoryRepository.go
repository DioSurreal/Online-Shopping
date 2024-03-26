package inventoryRepositories

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/inventory"
	itemPb "github.com/DioSurreal/Online-Shopping/modules/item/itemPb"
	"github.com/DioSurreal/Online-Shopping/modules/models"
	"github.com/DioSurreal/Online-Shopping/modules/payment"
	"github.com/DioSurreal/Online-Shopping/pkg/grpccon"
	"github.com/DioSurreal/Online-Shopping/pkg/jwtauth"
	"github.com/DioSurreal/Online-Shopping/pkg/queue"
	"github.com/DioSurreal/Online-Shopping/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	InventoryRepositoriesService interface {
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		FindUserItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventory.Inventory, error)
		CountUserItems(pctx context.Context, userId string) (int64, error)
		AddUserItemRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error
		RemoveUserItemRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error
		InsertOneUserItem(pctx context.Context, req *inventory.Inventory) (primitive.ObjectID, error)
		DeleteOneInventory(pctx context.Context, inventoryId string) error
		FindOneUserItem(pctx context.Context, userId, itemId string) bool
		DeleteOneUserItem(pctx context.Context, userId, itemId string) error
	}

	inventoryRepository struct {
		db *mongo.Client
	}
)

func NewInventoryRepository(db *mongo.Client) InventoryRepositoriesService {
	return &inventoryRepository{db}
}

func (r *inventoryRepository) inventoryDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("inventory")
}

func (r *inventoryRepository) GetOffset(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("users_inventory_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: GetOffset failed: %s", err.Error())
		return -1, errors.New("error: GetOffset failed")
	}

	return result.Offset, nil
}

func (r *inventoryRepository) UpserOffset(pctx context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("users_inventory_queue")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpserOffset failed: %s", err.Error())
		return errors.New("error: UpserOffset failed")
	}
	log.Printf("Info: UpserOffset result: %v", result)

	return nil
}

func (r *inventoryRepository) FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRPC connection failed: %s", err.Error())
		return nil, errors.New("error: gRPC connection failed")
	}

	result, err := conn.Item().FindItemsInIds(ctx, req)
	if err != nil {
		log.Printf("Error: FindItemsInIds failed: %s", err.Error())
		return nil, errors.New("error: items not found")
	}

	if result == nil {
		log.Printf("Error: FindItemsInIds failed: %s", err.Error())
		return nil, errors.New("error: items not found")
	}

	if result.Items == nil {
		log.Printf("Error: FindItemsInIds failed: %s", err.Error())
		return nil, errors.New("error: items not found")
	}

	if len(result.Items) == 0 {
		log.Printf("Error: FindItemsInIds failed: %s", err.Error())
		return nil, errors.New("error: items not found")
	}

	return result, nil
}

func (r *inventoryRepository) FindUserItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventory.Inventory, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("users_inventory")

	cursors, err := col.Find(ctx, filter, opts...)
	if err != nil {
		log.Printf("Error: FindUserItems failed: %s", err.Error())
		return nil, errors.New("error: user items not found")
	}

	results := make([]*inventory.Inventory, 0)
	for cursors.Next(ctx) {
		result := new(inventory.Inventory)
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: FindUserItems failed: %s", err.Error())
			return nil, errors.New("error: user items not found")
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *inventoryRepository) CountUserItems(pctx context.Context, userId string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("users_inventory")

	count, err := col.CountDocuments(ctx, bson.M{"user_id": userId})
	if err != nil {
		log.Printf("Error: CountUserItems failed: %s", err.Error())
		return -1, errors.New("error: count user items failed")
	}

	return count, nil
}

func (r *inventoryRepository) InsertOneUserItem(pctx context.Context, req *inventory.Inventory) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("users_inventory")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneUserItem failed: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert user item failed")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *inventoryRepository) DeleteOneInventory(pctx context.Context, inventoryId string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("users_inventory")

	result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(inventoryId)})
	if err != nil {
		log.Printf("Error: DeleteOneInventory failed: %s", err.Error())
		return errors.New("error: delete one inventory failed")
	}
	log.Printf("DeleteOneInventory result: %v", result)

	return nil
}

func (r *inventoryRepository) AddUserItemRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: AddUserItemRes failed: %s", err.Error())
		return errors.New("error: docked user money res failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"payment",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: AddUserItemRes failed: %s", err.Error())
		return errors.New("error: docked user money res failed")
	}

	return nil
}

func (r *inventoryRepository) FindOneUserItem(pctx context.Context, userId, itemId string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("users_inventory")

	result := new(inventory.Inventory)

	if err := col.FindOne(ctx, bson.M{"user_id": userId, "item_id": itemId}).Decode(result); err != nil {
		log.Printf("Error: FindOneUserItem failed: %s", err.Error())
		return false
	}
	return true
}

func (r *inventoryRepository) DeleteOneUserItem(pctx context.Context, userId, itemId string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("users_inventory")

	result, err := col.DeleteOne(ctx, bson.M{"user_id": userId, "item_id": itemId})
	if err != nil {
		log.Printf("Error: DeleteOneUserItem failed: %s", err.Error())
		return errors.New("error: delete one user item failed")
	}
	log.Printf("DeleteOneUserItem result: %v", result)

	return nil
}

func (r *inventoryRepository) RemoveUserItemRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: RemoveUserItemRes failed: %s", err.Error())
		return errors.New("error: docked user money res failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"payment",
		"sell",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RemoveUserItemRes failed: %s", err.Error())
		return errors.New("error: docked user money res failed")
	}

	return nil
}
