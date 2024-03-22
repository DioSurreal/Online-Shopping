package userUsecases

import "github.com/DioSurreal/Online-Shopping/modules/user/userRepositories"

type (
	UserUsecasesService interface{}

	userUsecase struct {
		userRepository userRepositories.UserRepositoriesService
	}
)

func NewUserUsecase (userRepository userRepositories.UserRepositoriesService) UserUsecasesService {
	return &userUsecase{userRepository: userRepository}
} 