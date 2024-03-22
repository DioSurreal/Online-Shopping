package userHandlers

import (
	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/user/userUsecases"
)
type (
	UserQueueHandlersService interface{}

	userQueueHandler struct {
		cfg  *config.Config
		userUsecase userUsecases.UserUsecasesService
	}
)

func NewUserQueueHandler(cfg *config.Config,userUsecase userUsecases.UserUsecasesService) UserQueueHandlersService {
	return &userQueueHandler{
		cfg: cfg,
		userUsecase: userUsecase,
	}
}