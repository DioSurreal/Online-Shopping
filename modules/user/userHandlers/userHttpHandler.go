package userHandlers

import (
	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/user/userUsecases"
)
type (
	UserHttpHandlersService interface{}

	userHttpHandler struct {
		cfg  *config.Config
		userUsecase userUsecases.UserUsecasesService
	}
)

func NewUserHttpHandler(cfg *config.Config,userUsecase userUsecases.UserUsecasesService) UserHttpHandlersService {
	return &userHttpHandler{
		cfg: cfg,
		userUsecase: userUsecase,
	}
}