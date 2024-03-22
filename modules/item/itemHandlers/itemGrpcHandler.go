package itemHandlers
import "github.com/DioSurreal/Online-Shopping/modules/item/itemUsecases"
type(
	itemGrpcHandler struct {
		itemUsecase itemUsecases.ItemUsecasesService
	}
)

func NewItemGrpcHandler(itemUsecase itemUsecases.ItemUsecasesService) itemGrpcHandler {
	return itemGrpcHandler{itemUsecase: itemUsecase}
}