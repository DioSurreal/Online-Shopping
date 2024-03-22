package userHandlers

import "github.com/DioSurreal/Online-Shopping/modules/user/userUsecases"
type(
	userGrpcHandler struct {
		userUsecase userUsecases.UserUsecasesService
	}
)

func NewUserGrpcHandler(userUsecase userUsecases.UserUsecasesService) userGrpcHandler {
	return userGrpcHandler{userUsecase: userUsecase}
}