package server

import (
	"github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryHandlers"
	"github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryRepositories"
	"github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryUsecases"
)

func (s *server) inventoryService() {
	inventoryRepository := inventoryRepositories.NewInventoryRepository(s.db)
	inventoryUsecase := inventoryUsecases.NewInventoryUsecase(inventoryRepository)
	inventoryHttpHandler := inventoryHandlers.NewInventoryHttpHandler(s.cfg,inventoryUsecase)
    inventoryGrpcHandler := inventoryHandlers.NewInventoryGrpcHandler(inventoryUsecase)

	_  = inventoryHttpHandler
	_ = inventoryGrpcHandler

	inventory := s.app.Group("/inventory_v1")

	//Health Check
	inventory.GET("",s.healthCheckService)
}