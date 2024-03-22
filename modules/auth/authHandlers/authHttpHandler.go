package authHandlers

import (
	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/auth/authUsecases"
)

type (
	AuthHandlersService interface{}

	authHttpHandlers struct {
		cfg         *config.Config
		authUsecase authUsecases.AuthUsecasesService
	}
)

func NewAuthHttpHandler(cfg *config.Config,authUsecase authUsecases.AuthUsecasesService) AuthHandlersService {
	return &authHttpHandlers{cfg,authUsecase}
}