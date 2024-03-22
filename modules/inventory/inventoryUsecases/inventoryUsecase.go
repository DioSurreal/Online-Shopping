package inventoryUsecases

import "github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryRepositories"

type(
	InventoryUsecasesService interface{}

	inventoryUsecase struct{
		inventoryRepository inventoryRepositories.InventoryRepositoriesService
	}
)

func NewInventoryUsecase (inventoryRepository inventoryRepositories.InventoryRepositoriesService) InventoryUsecasesService {
	return inventoryUsecase{inventoryRepository}
}

