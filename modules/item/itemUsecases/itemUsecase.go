package itemUsecases

import "github.com/DioSurreal/Online-Shopping/modules/item/itemRepositories"

type(
	ItemUsecasesService interface{}

	itemUsecase struct {
		itemRepository itemRepositories.ItemRepositoriesService
	}
)

func NewItemUsecase(itemRepository itemRepositories.ItemRepositoriesService) ItemUsecasesService {
	return &itemUsecase{itemRepository: itemRepository}
}