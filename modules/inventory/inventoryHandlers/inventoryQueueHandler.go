package inventoryHandlers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/inventory"
	"github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryUsecases"
	"github.com/DioSurreal/Online-Shopping/pkg/queue"
	"github.com/IBM/sarama"
)

type (
	InventoryQueueHandlerService interface {
		AddUserItem()
		RemoveUserItem()
		RollbackAddUserItem()
		RollbackRemoveUserItem()
	}

	inventoryQueueHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecases.InventoryUsecasesService
	}
)

func NewInventoryQueueHandler(cfg *config.Config, inventoryUsecase inventoryUsecases.InventoryUsecasesService) InventoryQueueHandlerService {
	return &inventoryQueueHandler{
		cfg:              cfg,
		inventoryUsecase: inventoryUsecase,
	}
}

func (h *inventoryQueueHandler) InventoryConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.inventoryUsecase.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("inventory", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("inventory", 0, 0)
		if err != nil {
			log.Println("Error: InventoryConsumer failed: ", err.Error())
			return nil, err
		}
	}

	return consumer, nil
}

func (h *inventoryQueueHandler) AddUserItem() {
	ctx := context.Background()

	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start AddUserItem ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: AddUserItem failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				h.inventoryUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(inventory.UpdateInventoryReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.AddUserItemRes(ctx, h.cfg, req)

				log.Printf("AddUserItem | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop AddUserItem...")
			return
		}
	}
}

func (h *inventoryQueueHandler) RollbackAddUserItem() {
	ctx := context.Background()

	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollbackAddUserItem ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackAddUserItem failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "radd" {
				h.inventoryUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(inventory.RollbackUserInventoryReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.RollbackAddUserItem(ctx, h.cfg, req)

				log.Printf("RollbackAddUserItem | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop RollbackAddUserItem...")
			return
		}
	}
}

func (h *inventoryQueueHandler) RemoveUserItem() {
	ctx := context.Background()

	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RemoveUserItem ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RemoveUserItem failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "sell" {
				h.inventoryUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(inventory.UpdateInventoryReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.RemoveUserItemRes(ctx, h.cfg, req)

				log.Printf("RemoveUserItem | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop RemoveUserItem...")
			return
		}
	}
}

func (h *inventoryQueueHandler) RollbackRemoveUserItem() {
	ctx := context.Background()

	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollbackRemoveUserItem ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackRemoveUserItem failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "rremove" {
				h.inventoryUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(inventory.RollbackUserInventoryReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.RollbackRemoveUserItem(ctx, h.cfg, req)

				log.Printf("RollbackRemoveUserItem | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop RollbackRemoveUserItem...")
			return
		}
	}
}