package server

import (
	"github.com/DioSurreal/Online-Shopping/modules/item/itemRepositories"
	"github.com/DioSurreal/Online-Shopping/modules/item/itemUsecases"
	"github.com/DioSurreal/Online-Shopping/modules/item/itemHandlers"
)

func (s *server) itemService() {
	itemRepository := itemRepositories.NewItemRepository(s.db)
	itemUsecase := itemUsecases.NewItemUsecase(itemRepository)
	itemHttpHandler := itemHandlers.NewItemHttpHandler(s.cfg,itemUsecase)
    itemGrpcHandler := itemHandlers.NewItemGrpcHandler(itemUsecase)

	_  = itemHttpHandler
	_ = itemGrpcHandler

	item := s.app.Group("/item_v1")

	//Health Check
	item.GET("",s.healthCheckService)
}