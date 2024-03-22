package authUsecases

import "github.com/DioSurreal/Online-Shopping/modules/auth/authRepositories"

type (
	AuthUsecasesService interface{}

	authUsecases struct {
		authRepository authRepositories.AuthRepositoryService
	}

)

func NewAuthUsecase(authRepository authRepositories.AuthRepositoryService) AuthUsecasesService {
	return & authUsecases{authRepository}
}