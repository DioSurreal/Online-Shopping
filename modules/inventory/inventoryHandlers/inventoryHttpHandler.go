package inventoryHandlers

import (
	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryUsecases"
)

type(
	InventoryHttpHandlersService interface{}

	inventoryHttpHandler struct{
		cfg  *config.Config
		inventoryUsecase inventoryUsecases.InventoryUsecasesService
	}
)

func NewInventoryHttpHandler (cfg *config.Config,inventoryUsecase inventoryUsecases.InventoryUsecasesService) InventoryHttpHandlersService {
	return inventoryHttpHandler{
		cfg: cfg,
		inventoryUsecase: inventoryUsecase,}
}