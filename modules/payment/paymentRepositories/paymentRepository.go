package paymentRepositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type(
	PaymentRepositoriesService interface{}

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