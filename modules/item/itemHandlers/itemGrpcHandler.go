package itemHandlers

import (
	"context"

	itemPb "github.com/DioSurreal/Online-Shopping/modules/item/itemPb"
	"github.com/DioSurreal/Online-Shopping/modules/item/itemUsecases"
)
type(
	itemGrpcHandler struct {
		itemUsecase itemUsecases.ItemUsecasesService
		itemPb.UnimplementedItemGrpcServiceServer
	}
)

func NewItemGrpcHandler(itemUsecase itemUsecases.ItemUsecasesService) *itemGrpcHandler {
	return &itemGrpcHandler{itemUsecase: itemUsecase}
}

func (g *itemGrpcHandler) FindItemsInIds(ctx context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	return nil,nil
}