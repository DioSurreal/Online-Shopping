package paymentRepositories

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
	"github.com/DioSurreal/Online-Shopping/modules/user"
	"github.com/DioSurreal/Online-Shopping/pkg/grpccon"
	"github.com/DioSurreal/Online-Shopping/pkg/jwtauth"
	"github.com/DioSurreal/Online-Shopping/pkg/queue"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type(
	PaymentRepositoriesService interface{
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		DockedUserMoney(pctx context.Context, cfg *config.Config, req *user.CreateUserTransactionReq) error
		RollbackTransaction(pctx context.Context, cfg *config.Config, req *user.RollbackUserTransactionReq) error
		AddUserItem(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) error
		RollbackAddUserItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackUserInventoryReq) error
		RemoveUserItem(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) error
		RollbackRemoveUserItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackUserInventoryReq) error
		AddUserMoney(pctx context.Context, cfg *config.Config, req *user.CreateUserTransactionReq) error
	}

    paymentRepository struct {
		db *mongo.Client
	}
)

func NewPaymentRepository (db *mongo.Client) PaymentRepositoriesService{
	return &paymentRepository{db}
}

func (r *paymentRepository) paymentDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("payment_db")
}
func (r *paymentRepository) GetOffset(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payment_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: GetOffset failed: %s", err.Error())
		return -1, errors.New("error: GetOffset failed")
	}

	return result.Offset, nil
}

func (r *paymentRepository) UpserOffset(pctx context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payment_queue")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpserOffset failed: %s", err.Error())
		return errors.New("error: UpserOffset failed")
	}
	log.Printf("Info: UpserOffset result: %v", result)

	return nil
}

func (r *paymentRepository) FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
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

	if len(result.Items) == 0 {
		log.Printf("Error: FindItemsInIds failed: %s", err.Error())
		return nil, errors.New("error: items not found")
	}

	return result, nil
}

func (r *paymentRepository) DockedUserMoney(pctx context.Context, cfg *config.Config, req *user.CreateUserTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: DockedUserMoney failed: %s", err.Error())
		return errors.New("error: docked user money failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"user",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: DockedUserMoney failed: %s", err.Error())
		return errors.New("error: docked user money failed")
	}

	return nil
}

func (r *paymentRepository) AddUserMoney(pctx context.Context, cfg *config.Config, req *user.CreateUserTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: AddUserMoney failed: %s", err.Error())
		return errors.New("error: add user money failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"user",
		"sell",
		reqInBytes,
	); err != nil {
		log.Printf("Error: AddUserMoney failed: %s", err.Error())
		return errors.New("error: add user money failed")
	}

	return nil
}

func (r *paymentRepository) RollbackTransaction(pctx context.Context, cfg *config.Config, req *user.RollbackUserTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: DockedUserMoney failed: %s", err.Error())
		return errors.New("error: rollback user transaction failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"user",
		"rtransaction",
		reqInBytes,
	); err != nil {
		log.Printf("Error: DockedUserMoney failed: %s", err.Error())
		return errors.New("error: rollback user transaction failed")
	}

	return nil
}

func (r *paymentRepository) AddUserItem(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: AddUserItem failed: %s", err.Error())
		return errors.New("error: add user item failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: AddUserItem failed: %s", err.Error())
		return errors.New("error: add user item failed")
	}

	return nil
}

func (r *paymentRepository) RollbackAddUserItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackUserInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: RollbackAddUserItem failed: %s", err.Error())
		return errors.New("error: rollback add user item failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"radd",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RollbackAddUserItem failed: %s", err.Error())
		return errors.New("error: rollback add user item failed")
	}

	return nil
}

func (r *paymentRepository) RemoveUserItem(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: RemoveUserItem failed: %s", err.Error())
		return errors.New("error: remove user item failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"sell",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RemoveUserItem failed: %s", err.Error())
		return errors.New("error: remove user item failed")
	}

	return nil
}

func (r *paymentRepository) RollbackRemoveUserItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackUserInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: RollbackRemoveUserItem failed: %s", err.Error())
		return errors.New("error: rollback remove user item failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"rremove",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RollbackRemoveUserItem failed: %s", err.Error())
		return errors.New("error: rollback remove user item failed")
	}

	return nil
}