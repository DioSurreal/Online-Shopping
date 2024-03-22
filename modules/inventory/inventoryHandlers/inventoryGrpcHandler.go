package inventoryHandlers

import "github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryUsecases"

type(
	inventoryGrpcHandler struct{
		inventoryUsecase inventoryUsecases.InventoryUsecasesService
	}
)

func NewInventoryGrpcHandler (inventoryUsecase inventoryUsecases.InventoryUsecasesService) inventoryGrpcHandler {
	return inventoryGrpcHandler{inventoryUsecase}
}