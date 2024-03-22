package itemHandlers

import (
	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/item/itemUsecases"
)
type(
	ItemHttpHandlersService interface{}

	itemHttpHandler struct {
		cfg  *config.Config
		itemUsecase itemUsecases.ItemUsecasesService
	}
)

func NewItemHttpHandler(cfg *config.Config,itemUsecase itemUsecases.ItemUsecasesService) ItemHttpHandlersService {
	return &itemHttpHandler{
		cfg: cfg,
		itemUsecase: itemUsecase}
}