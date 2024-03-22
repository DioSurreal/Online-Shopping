package middlewareHandlers

import (
	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/middleware/middlewareUsecases"
)

type (
	MiddlewareHandlersService interface{}

	middlewareHandler struct {
		cfg               *config.Config
		middlewareUsecase middlewareUsecases.MiddlewareUsecasesService
	}
)

func NewMiddlewareHandler(cfg *config.Config, middlewareUsecase middlewareUsecases.MiddlewareUsecasesService) MiddlewareHandlersService {
	return &middlewareHandler{
		cfg:               cfg,
		middlewareUsecase: middlewareUsecase,}
}
