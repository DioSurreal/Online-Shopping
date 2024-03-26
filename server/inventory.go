package server

import (
	"github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryHandlers"
	"github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryRepositories"
	"github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryUsecases"
)

func (s *server) inventoryService() {
	inventoryRepository := inventoryRepositories.NewInventoryRepository(s.db)
	inventoryUsecase := inventoryUsecases.NewInventoryUsecase(inventoryRepository)
	inventoryHttpHandler := inventoryHandlers.NewInventoryHttpHandler(s.cfg, inventoryUsecase)
	inventoryqueueHandler := inventoryHandlers.NewInventoryQueueHandler(s.cfg, inventoryUsecase)

	go inventoryqueueHandler.AddUserItem()
	go inventoryqueueHandler.RollbackAddUserItem()
	go inventoryqueueHandler.RemoveUserItem()
	go inventoryqueueHandler.RollbackRemoveUserItem()

	inventory := s.app.Group("/inventory_v1")

	//Health Check
	inventory.GET("", s.healthCheckService)
	inventory.GET("/inventory/:user_id", inventoryHttpHandler.FindUserItems, s.middleware.JwtAuthorization, s.middleware.UserIdParamValidation)
}
