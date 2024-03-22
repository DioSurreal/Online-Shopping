package middlewareUsecases

import middlewareRepository "github.com/DioSurreal/Online-Shopping/modules/middleware/middlewareRepositories"

type (
	MiddlewareUsecasesService interface{}

	middlewareUsecase struct{
		middlewareRepository middlewareRepository.MiddlewareRepositoriesService
	}
)

func NewMiddlewareUsecase(middlewareRepository middlewareRepository.MiddlewareRepositoriesService) MiddlewareUsecasesService {
	return &middlewareUsecase{middlewareRepository}
}