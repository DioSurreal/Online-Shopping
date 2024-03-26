package inventoryUsecases

import (
	"context"
	"fmt"

	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/inventory"
	"github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryRepositories"
	"github.com/DioSurreal/Online-Shopping/modules/item"
	itemPb "github.com/DioSurreal/Online-Shopping/modules/item/itemPb"
	"github.com/DioSurreal/Online-Shopping/modules/models"
	"github.com/DioSurreal/Online-Shopping/modules/payment"
	"github.com/DioSurreal/Online-Shopping/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	InventoryUsecasesService interface {
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		FindUserItems(pctx context.Context, cfg *config.Config, userId string, req *inventory.InventorySearchReq) (*models.PaginateRes, error)
		AddUserItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq)
		RemoveUserItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq)
		RollbackAddUserItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackUserInventoryReq)
		RollbackRemoveUserItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackUserInventoryReq)
	}

	inventoryUsecase struct {
		inventoryRepository inventoryRepositories.InventoryRepositoriesService
	}
)

func NewInventoryUsecase(inventoryRepository inventoryRepositories.InventoryRepositoriesService) InventoryUsecasesService {
	return &inventoryUsecase{inventoryRepository}
}

func (u *inventoryUsecase) GetOffset(pctx context.Context) (int64, error) {
	return u.inventoryRepository.GetOffset(pctx)
}
func (u *inventoryUsecase) UpserOffset(pctx context.Context, offset int64) error {
	return u.inventoryRepository.UpserOffset(pctx, offset)
}

func (u *inventoryUsecase) FindUserItems(pctx context.Context, cfg *config.Config, userId string, req *inventory.InventorySearchReq) (*models.PaginateRes, error) {
	// Filter
	filter := bson.D{}

	// Filter
	if req.Start != "" {
		filter = append(filter, bson.E{"_id", bson.D{{"$gt", utils.ConvertToObjectId(req.Start)}}})
	}
	filter = append(filter, bson.E{"user_id", userId})

	// Option
	opts := make([]*options.FindOptions, 0)

	opts = append(opts, options.Find().SetSort(bson.D{{"_id", 1}}))
	opts = append(opts, options.Find().SetLimit(int64(req.Limit)))

	// Find
	inventoryData, err := u.inventoryRepository.FindUserItems(pctx, filter, opts)
	if err != nil {
		return nil, err
	}
	if len(inventoryData) == 0 {
		return &models.PaginateRes{
			Data:  make([]*inventory.ItemInInventory, 0),
			Total: 0,
			Limit: req.Limit,
			First: models.FirstPaginate{
				Href: fmt.Sprintf("%s/%s?limit=%d", cfg.Paginate.InventoryNextPageBasedUrl, userId, req.Limit),
			},
			Next: models.NextPaginate{
				Start: "",
				Href:  "",
			},
		}, nil
	}

	itemData, err := u.inventoryRepository.FindItemsInIds(pctx, cfg.Grpc.ItemUrl, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0)
			for _, v := range inventoryData {
				itemIds = append(itemIds, v.ItemId)
			}
			return itemIds
		}(),
	})

	itemMaps := make(map[string]*item.ItemShowCase)
	for _, v := range itemData.Items {
		itemMaps[v.Id] = &item.ItemShowCase{
			ItemId:   v.Id,
			Title:    v.Title,
			Price:    v.Price,
			ImageUrl: v.ImageUrl,
			Damage:   int(v.Damage),
		}
	}

	results := make([]*inventory.ItemInInventory, 0)
	for _, v := range inventoryData {
		results = append(results, &inventory.ItemInInventory{
			InventoryId: v.Id.Hex(),
			UserId:      v.UserId,
			ItemShowCase: &item.ItemShowCase{
				ItemId:   v.ItemId,
				Title:    itemMaps[v.ItemId].Title,
				Price:    itemMaps[v.ItemId].Price,
				Damage:   itemMaps[v.ItemId].Damage,
				ImageUrl: itemMaps[v.ItemId].ImageUrl,
			},
		})
	}

	// Count
	total, err := u.inventoryRepository.CountUserItems(pctx, userId)
	if err != nil {
		return nil, err
	}

	return &models.PaginateRes{
		Data:  results,
		Total: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s/%s?limit=%d", cfg.Paginate.InventoryNextPageBasedUrl, userId, req.Limit),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].InventoryId,
			Href:  fmt.Sprintf("%s/%s?limit=%d&start=%s", cfg.Paginate.InventoryNextPageBasedUrl, userId, req.Limit, results[len(results)-1].InventoryId),
		},
	}, nil
}

func (u *inventoryUsecase) AddUserItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) {
	inventoryId, err := u.inventoryRepository.InsertOneUserItem(pctx, &inventory.Inventory{
		UserId: req.UserId,
		ItemId: req.ItemId,
	})
	if err != nil {
		u.inventoryRepository.AddUserItemRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			UserId:        req.UserId,
			ItemId:        req.ItemId,
			Amount:        0,
			Error:         err.Error(),
		})
		return
	}

	u.inventoryRepository.AddUserItemRes(pctx, cfg, &payment.PaymentTransferRes{
		InventoryId:   inventoryId.Hex(),
		TransactionId: "",
		UserId:        req.UserId,
		ItemId:        req.ItemId,
		Amount:        0,
		Error:         "",
	})
}

func (u *inventoryUsecase) RemoveUserItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) {
	if !u.inventoryRepository.FindOneUserItem(pctx, req.UserId, req.ItemId) {
		u.inventoryRepository.RemoveUserItemRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			UserId:        req.UserId,
			ItemId:        req.ItemId,
			Amount:        0,
			Error:         "error: item not found",
		})
		return
	}

	if err := u.inventoryRepository.DeleteOneUserItem(pctx, req.UserId, req.ItemId); err != nil {
		u.inventoryRepository.RemoveUserItemRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			UserId:        req.UserId,
			ItemId:        req.ItemId,
			Amount:        0,
			Error:         err.Error(),
		})
		return
	}

	u.inventoryRepository.RemoveUserItemRes(pctx, cfg, &payment.PaymentTransferRes{
		InventoryId:   "",
		TransactionId: "",
		UserId:        req.UserId,
		ItemId:        req.ItemId,
		Amount:        0,
		Error:         "",
	})
}

func (u *inventoryUsecase) RollbackAddUserItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackUserInventoryReq) {
	u.inventoryRepository.DeleteOneInventory(pctx, req.InventoryId)
}

func (u *inventoryUsecase) RollbackRemoveUserItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackUserInventoryReq) {
	u.inventoryRepository.InsertOneUserItem(pctx, &inventory.Inventory{
		UserId: req.UserId,
		ItemId: req.ItemId,
	})
}
